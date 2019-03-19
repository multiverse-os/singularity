package memfd

import (
	"errors"
	"os"
	"syscall"

	"github.com/multiverse-os/singularity/msyscall"
)

var (
	ErrTooBig = errors.New("memfd too large for slice")
)

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

func (mfd *MemFD) Size() int64 {
	fi, err := mfd.Stat()
	if err != nil {
		return 0
	}
	return fi.Size()
}

func (mfd *MemFD) SetSize(size int64) error {
	return mfd.Truncate(size)
}

func (mfd *MemFD) ClearCloexec() {
	_ = msyscall.FcntlCloexec(mfd.Fd(), 0)
}

func (mfd *MemFD) SetCloexec() {
	_ = msyscall.FcntlCloexec(mfd.Fd(), 1)
}

func (mfd *MemFD) seals() (int, error) {
	return msyscall.FcntlSeals(mfd.Fd())
}

func (mfd *MemFD) Seals() int {
	seals, err := mfd.seals()
	if err != nil {
		return 0
	}
	return seals
}

func (mfd *MemFD) SetSeals(seals int) error {
	return msyscall.FcntlSetSeals(mfd.Fd(), seals)
}

func (mfd *MemFD) IsImmutable() bool {
	seals, err := msyscall.FcntlSeals(mfd.Fd())
	if err != nil {
		return false
	}
	return seals == SealAll
}

func (mfd *MemFD) SetImmutable() error {
	err := mfd.SetSeals(SealAll)
	if err == nil {
		return nil
	}
	if mfd.IsImmutable() {
		return nil
	}
	return err
}

func (mfd *MemFD) Map() ([]byte, error) {
	if cap(mfd.Bytes) > 0 {
		return mfd.Bytes, nil
	}
	seals, err := mfd.seals()
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
	size := mfd.Size()
	if size > maxint {
		return []byte{}, ErrTooBig
	}
	if size == 0 {
		return []byte{}, nil
	}
	bytes, err := syscall.Mmap(int(mfd.Fd()), 0, int(size), prot, flags)
	if err != nil {
		return []byte{}, err
	}
	mfd.Bytes = bytes
	return bytes, nil
}

func (mfd *MemFD) Unmap() error {
	if cap(mfd.Bytes) == 0 {
		return nil
	}
	err := syscall.Munmap(mfd.Bytes)
	mfd.Bytes = []byte{}
	return err
}

func (mfd *MemFD) Remap() ([]byte, error) {
	if cap(mfd.Bytes) == 0 {
		return mfd.Map()
	}
	err := mfd.Unmap()
	if err != nil {
		return []byte{}, err
	}
	return mfd.Map()
}
