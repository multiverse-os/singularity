package memfd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/multiverse-os/singularity/msyscall"
)

var ErrTooBig = errors.New("memfd too large for slice")

const maxint int64 = int64(^uint(0) >> 1)

const (
	Cloexec      = msyscall.MFD_CLOEXEC
	AllowSealing = msyscall.MFD_ALLOW_SEALING
	SealSeal     = msyscall.F_SEAL_SEAL
	SealShrink   = msyscall.F_SEAL_SHRINK
	SealGrow     = msyscall.F_SEAL_GROW
	SealWrite    = msyscall.F_SEAL_WRITE
	SealAll      = SealSeal | SealShrink | SealGrow | SealWrite
)

type MemFD struct {
	*os.File
	Bytes []byte
}

func (self *MemFD) FDPath() string {
	return fmt.Sprintf("/proc/self/fd/%d", int(self.Fd()))
}

func (self *MemFD) FDFileInfo() (os.FileInfo, error) {
	return os.Lstat(self.FDPath())
}

func (self *MemFD) FDPerm() (os.FileMode, error) {
	fi, err := self.FDFileInfo()
	return fi.Mode().Perm(), err
}

func (self *MemFD) Command(arg ...string) *exec.Cmd {
	return exec.Command(self.FDPath(), arg...)
}

func (self *MemFD) Close() error {
	return self.File.Close()
}

func Create() (*MemFD, error) {
	return CreateNameFlags("", Cloexec|AllowSealing)
}

func CreateNameFlags(name string, flags uint) (*MemFD, error) {
	fd, err := msyscall.MemFDCreate(name, flags)
	if err != nil {
		return nil, err
	}
	memfd := MemFD{os.NewFile(uintptr(fd), name), []byte{}}
	return &memfd, nil
}

func New(fd uintptr) (*MemFD, error) {
	_, err := msyscall.FcntlSeals(fd)
	if err != nil {
		return nil, err
	}
	// TODO(justin) read name with readlink /proc/self/fd
	mfd := MemFD{os.NewFile(uintptr(fd), ""), []byte{}}
	return &mfd, nil
}

func (self *MemFD) Size() int64 {
	fi, err := self.Stat()
	if err != nil {
		return 0
	}
	return fi.Size()
}

func (self *MemFD) SetSize(size int64) error {
	return self.Truncate(size)
}

func (self *MemFD) ClearCloexec() {
	_ = msyscall.FcntlCloexec(self.Fd(), 0)
}

func (self *MemFD) SetCloexec() {
	_ = msyscall.FcntlCloexec(self.Fd(), 1)
}

func (self *MemFD) seals() (int, error) {
	return msyscall.FcntlSeals(self.Fd())
}

func (self *MemFD) Seals() int {
	seals, err := self.seals()
	if err != nil {
		return 0
	}
	return seals
}

func (self *MemFD) SetSeals(seals int) error {
	return msyscall.FcntlSetSeals(self.Fd(), seals)
}

func (self *MemFD) IsImmutable() bool {
	seals, err := msyscall.FcntlSeals(self.Fd())
	if err != nil {
		return false
	}
	return seals == SealAll
}

func (self *MemFD) SetImmutable() error {
	err := self.SetSeals(SealAll)
	if err == nil {
		return nil
	}
	if self.IsImmutable() {
		return nil
	}
	return err
}

func (self *MemFD) Map() ([]byte, error) {
	if cap(self.Bytes) > 0 {
		return self.Bytes, nil
	}
	seals, err := self.seals()
	if err != nil {
		return []byte{}, err
	}
	var prot, flags int
	if seals&SealWrite == SealWrite {
		prot = syscall.PROT_READ
		flags = syscall.MAP_PRIVATE
	} else {
		prot = syscall.PROT_READ | syscall.PROT_WRITE
		flags = syscall.MAP_SHARED
	}
	size := self.Size()
	if size > maxint {
		return []byte{}, ErrTooBig
	}
	if size == 0 {
		return []byte{}, nil
	}
	bytes, err := syscall.Mmap(int(self.Fd()), 0, int(size), prot, flags)
	if err != nil {
		return []byte{}, err
	}
	self.Bytes = bytes
	return bytes, nil
}

func (self *MemFD) Unmap() error {
	if cap(self.Bytes) == 0 {
		return nil
	}
	err := syscall.Munmap(self.Bytes)
	self.Bytes = []byte{}
	return err
}

func (self *MemFD) Remap() ([]byte, error) {
	if cap(self.Bytes) == 0 {
		return self.Map()
	}
	err := self.Unmap()
	if err != nil {
		return []byte{}, err
	}
	return self.Map()
}
