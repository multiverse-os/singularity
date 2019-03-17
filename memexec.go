package memexec

import (
	"io/ioutil"
	"os"
	"os/exec"
)

type mem struct {
	f *os.File
}

func New(b []byte) (*mem, error) {
	f, err := ioutil.TempFile("", "go-memexec-")
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if f != nil && err != nil {
			f.Close()
			os.Remove(f.Name())
		}
	}(f)
	if err = os.Chmod(f.Name(), 0500); err != nil {
		return nil, err
	}
	if f, err = write(f, b); err != nil {
		return nil, err
	}
	return &mem{f: f}, nil
}

func (m *mem) Command(args ...string) *exec.Cmd {
	return exec.Command(path(m), args...)
}

func (m *mem) Close() error {
	return close(m)
}
