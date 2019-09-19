package prefixfs

import (
	"os"

	vfs "github.com/multiverse-os/singularity/store/vfs"
)

// Prefix is used to prefix the path in each vfs.Filesystem operation.
type FS struct {
	vfs.Filesystem
	Prefix string
}

func Create(root vfs.Filesystem, prefix string) *FS { return &FS{root, prefix} }
func (fs *FS) PrefixPath(path string) string        { return fs.Prefix + string(fs.PathSeparator()) + path }
func (fs *FS) PathSeparator() uint8                 { return fs.Filesystem.PathSeparator() }

func (fs *FS) OpenFile(name string, flag int, perm os.FileMode) (vfs.File, error) {
	return fs.Filesystem.OpenFile(fs.PrefixPath(name), flag, perm)
}

func (fs *FS) Remove(name string) error { return fs.Filesystem.Remove(fs.PrefixPath(name)) }

func (fs *FS) Rename(oldpath, newpath string) error {
	return fs.Filesystem.Rename(fs.PrefixPath(oldpath), fs.PrefixPath(newpath))
}

func (fs *FS) Mkdir(name string, perm os.FileMode) error {
	return fs.Filesystem.Mkdir(fs.PrefixPath(name), perm)
}

func (fs *FS) Stat(name string) (os.FileInfo, error)  { return fs.Filesystem.Stat(fs.PrefixPath(name)) }
func (fs *FS) Lstat(name string) (os.FileInfo, error) { return fs.Filesystem.Lstat(fs.PrefixPath(name)) }

func (fs *FS) ReadDir(path string) ([]os.FileInfo, error) {
	return fs.Filesystem.ReadDir(fs.PrefixPath(path))
}
