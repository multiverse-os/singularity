package singularity

import (
	"fmt"
	"os"
	"os/exec"

	memfd "github.com/multiverse-os/singularity/memfd"
	memfs "github.com/multiverse-os/singularity/memfs"
	msyscall "github.com/multiverse-os/singularity/msyscall"
)

var CurrentProcess Process

type Process struct {
	PID      int
	Binaries map[string]*Binary
}

type Binary struct {
	ParentProcess *Process
	Path          string
	Name          string
	Size          int
	Data          []byte
	Output        string
	ExitCode      int
	Permissions   os.FileMode
	MemFD         memfd.MemFD
	MemFS         *memfs.Filesystem
}

func LoadBinary(binaryName string, data []byte) *Binary {
	fs, err := memfs.NewFS()
	if err != nil {
		fmt.Println("[fatal] failed to open memory filesystem:", err)
		os.Exit(1)
	}
	CurrentProcess := &Process{
		PID:      os.Getpid(),
		Binaries: make(map[string]*Binary),
	}
	binary := &Binary{
		ParentProcess: CurrentProcess,
		Name:          binaryName,
		Size:          len(data),
		Data:          data,
		MemFS:         fs,
	}
	CurrentProcess.Binaries[binary.Name] = binary
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
	mfdBytes, err := self.MemFD.Map()
	if err != nil || self.Size != len(mfdBytes) {
		fmt.Println("[fatal error] failed to create memory fd or write binary data to fd:", err)
		os.Exit(1)
	}
	perm, err := self.MemFD.FDPerm()
	if err != nil {
		fmt.Println("[fatal error] failed to obtain mem FD permissions:", err)
		os.Exit(1)
	} else {
		self.Permissions = perm
	}

	fmt.Println("[singularity] mapped binary fd, bytes written to fd: [", len(mfdBytes), "b ]")
	fmt.Println("[singularity] fd path:", self.MemFD.FDPath())
	self.PrintString()

	return self
}

func (self *Binary) Execute(arguments ...string) ([]byte, error) {
	fmt.Println("[singularity] Inside Execute()")
	//self.MemFD.SetCloexec()
	return exec.Command(self.MemFD.Name(), arguments...).Output()
}

func (self *Binary) PrintString() {
	fmt.Println("===[LOADED MemFD Executable ]================")
	fmt.Println("| |-[Name] ", self.Name)
	fmt.Println("| |-[Size] ", self.Size)
	fmt.Println("=============================================")
}
