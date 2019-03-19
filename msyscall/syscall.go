package msyscall

import (
	"os"
	"syscall"
	"unsafe"
)

const (
	MFD_CLOEXEC       uint = 1
	MFD_ALLOW_SEALING uint = 2

	f_LINUX_SPECIFIC_BASE int = 1024

	F_ADD_SEALS int = (f_LINUX_SPECIFIC_BASE + 9)
	F_GET_SEALS int = (f_LINUX_SPECIFIC_BASE + 10)

	F_SEAL_SEAL   int = 0x0001
	F_SEAL_SHRINK int = 0x0002
	F_SEAL_GROW   int = 0x0004
	F_SEAL_WRITE  int = 0x0008
)

func MemFDCreate(name string, flags uint) (fd uintptr, err error) {
	var _p0 *byte
	_p0, err = syscall.BytePtrFromString(name)
	if err != nil {
		return
	}
	fd, _, e1 := syscall.Syscall(SYS_MEMFD_CREATE, uintptr(unsafe.Pointer(_p0)), uintptr(flags), uintptr(0))
	if e1 != 0 {
		err = os.NewSyscallError("memfd_create", e1)
	}
	return
}

func FcntlSeals(fd uintptr) (seals int, err error) {
	r1, _, e1 := syscall.Syscall(SYS_FCNTL, fd, uintptr(F_GET_SEALS), uintptr(0))
	if e1 != 0 {
		err = os.NewSyscallError("fcntl", e1)
	}
	seals = int(r1)
	return
}

func FcntlSetSeals(fd uintptr, seals int) (err error) {
	_, _, e1 := syscall.Syscall(SYS_FCNTL, fd, uintptr(F_ADD_SEALS), uintptr(seals))
	if e1 != 0 {
		err = os.NewSyscallError("fcntl", e1)
	}
	return
}

func FcntlCloexec(fd uintptr, flag int) (err error) {
	_, _, e1 := syscall.Syscall(SYS_FCNTL, fd, syscall.F_SETFD, uintptr(flag))
	if e1 != 0 {
		err = os.NewSyscallError("fcntl", e1)
	}
	return
}
