package tmpfs

import (
	"errors"
	"syscall"
)

const TMPFS_MAGIC = 0x01021994

// TODO: Add some help functions to make the size aspect easier to work with
func Mount(path string, size int64) error {
	if size < 0 {
		return errors.New("[tmpfs] tmpfs.Mount(path, size): size < 0")
	}
	var flags uintptr
	flags = syscall.MS_NOATIME | syscall.MS_SILENT
	// TODO: Make NOEXEC optional
	flags |= syscall.MS_NOSUID
	options = "size=" + strconv.FormatInt(size, 10)
	err := syscall.Mount("tmpfs", path, "tmpfs", flags, options)
	return os.NewSyscallError("mount", err)
}

func Unmount(path string) error { return syscall.Unmount(path, 0) }
