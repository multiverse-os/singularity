package singularity

import (
	"fmt"
	"os"
	"os/exec"

	memfd "github.com/multiverse-os/singularity/memfd"
	memfs "github.com/multiverse-os/singularity/memfs"
	msyscall "github.com/multiverse-os/singularity/msyscall"
)

// TODO: Singularity Tasks:
// 1) Take the test.rb stored as binary, and put it in the MemFS and use the ruby
// binary executed solely in memory to run the script in memFS
// 1a) Write exit code to binary
// 2a) Test long running processes and manage them, relaunch, etc

// 2) Add files to binary and fill them with data during runtime
// 3) Run ruby scripts from loaded file abstracts loaded during runtime (similar
// to patching, will be used for patching and other tasks in Multiverse OS)
// 4) Add compresssiona and cryptography middleware to data

var CurrentProcess Process

type Process struct {
	PID      int
	Binaries map[string]*Binary
}

type Binary struct {
	ParentProcess *Process
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
		fmt.Println("[fatal] failed to create memory fd or write binary data to fd:", err)
		os.Exit(1)
	}
	perm, err := self.MemFD.FDPerm()
	if err != nil {
		fmt.Println("[fatal] failed to obtain mem FD permissions:", err)
		os.Exit(1)
	} else {
		self.Permissions = perm
	}
	fmt.Println("[singularity] mapped binary fd, bytes written to fd: [", len(mfdBytes), "b ]")
	self.PrintString()
	return self
}

func (self *Binary) Execute(arguments ...string) ([]byte, error) {
	fmt.Println("[singularity] Inside Execute()")
	return exec.Command(self.MemFD.Name(), arguments...).Output()
}

func (self *Binary) PrintString() {
	fmt.Println("===[LOADED MemFD Executable ]================")
	fmt.Println("| |-[Name   ] ", self.Name)
	fmt.Println("| |-[Size   ] ", self.Size)
	fmt.Println("| |-[FD Path] ", self.MemFD.FDPath())
	fmt.Println("=============================================")
}
