package singularity

import (
	"fmt"
	"os"
	//"os/exec"

	memfd "github.com/multiverse-os/singularity/memfd"
	msyscall "github.com/multiverse-os/singularity/msyscall"
)

type Binary struct {
	Name     string
	Size     int
	Data     []byte
	Output   string
	ExitCode int
	MemFD    memfd.MemFD
}

func (self Binary) FDPath() string {
	return fmt.Sprintf("/proc/self/fd/%d", int(self.MemFD.Fd()))
}

func LoadBinary(name string, data []byte) *Binary {
	binary := &Binary{
		Name: name,
		Size: len(data),
		Data: data,
	}
	binary = binary.NewMemFD()
	fmt.Println("[singularity] loaded binary: [", binary.Name, "] with size of [", binary.Size, "bytes ]")
	return binary
}

func (self *Binary) NewMemFD() *Binary {
	fd, err := msyscall.MemFDCreate(self.Name, memfd.Cloexec|memfd.AllowSealing)
	if err != nil {
		fmt.Println("[error] failed to create memory fd:", err)
	}
	self.MemFD = memfd.MemFD{os.NewFile(uintptr(fd), self.Name), self.Data}
	//self.MemFD = MemFD{os.NewFile(uintptr(fd), self.Name)}
	mfdBytes, err := self.MemFD.Map()
	if err != nil || self.Size != len(mfdBytes) {
		fmt.Println("[fatal error] failed to create memory fd or write binary data to fd:", err)
		os.Exit(1)
	}
	fmt.Println("[singularity] mapped binary fd, bytes written to fd: [", len(mfdBytes), "b ]")
	fmt.Println("[singularity] fd path:", self.FDPath())
	return self
}

func (self *Binary) Execute() (string, error) {
	fmt.Println("[singularity] Inside Execute()")
	return "", nil
}
