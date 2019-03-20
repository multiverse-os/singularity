package singularity

// NOTE: ALternative way to update,
// store the soruce code, use it to rebuild and replace itself

// TODO:
// * Add JSON output of all files with bool to include data
// * Ls() to list files
// * Add actions for files to files as methods
// * Add locking to files so its threadsafe
// * Add in memory pub/sub for live updates on file changes for all subscribes
// * Have a map in the FS that can pull any file out by either their checksum or
// something similar.
// * Automatically store a checksum for each file
// * Add IsExecutable()
// * Add mime/magic based file type guesing
// * Cat() ability to pipe
// * Tail()
// * AppendToFile()
// * Add Path()

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

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
	ParentProcess  *Process
	ExecutablePath string
	Path           string
	Name           string
	Size           int
	Data           []byte
	Output         string
	ExitCode       int
	Permissions    os.FileMode
	MemFD          memfd.MemFD
	MemFS          *memfs.Filesystem
}

func LoadBinary(p string, data []byte) *Binary {
	fs, err := memfs.NewFS()
	if err != nil {
		fmt.Println("[fatal error] failed to open memory filesystem:", err)
		os.Exit(1)
	}
	CurrentProcess := &Process{
		PID:      os.Getpid(),
		Binaries: make(map[string]*Binary),
	}
	path, filename := filepath.Split(p)
	binary := &Binary{
		ParentProcess:  CurrentProcess,
		ExecutablePath: p,
		Path:           path,
		Name:           filename,
		Size:           len(data),
		Data:           data,
		MemFS:          fs,
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
	//self.MemFD = MemFD{os.NewFile(uintptr(fd), self.Name)}
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

func (self *Binary) Execute() (string, error) {
	fmt.Println("[singularity] Inside Execute()")
	//filepath := fmt.Sprintf("/proc/%d/fd/%d", pid, self.MemFD)
	//f, err := os.OpenFile(self.MemFD.FDPath(), os.O_RDWR, 0755)
	//if err != nil {
	//	fmt.Println("[error] failed to obtain permissions for memory fd:", err)
	//}
	//fmt.Println("fd permissions are:", perm.String())
	//self.MemFD.SetCloexec()

	// Opens a file with read/write permissions in the current directory
	fmt.Println("self.Name(): ", self.Name)
	f, err := self.MemFS.Create("ruby")
	if err != nil {
		fmt.Println("[error] failed to create memfs file:", err)
	}

	f.Write(self.Data)
	defer f.Close()

	fmt.Println("f data:", f.Name())
	cmd := exec.Command(f.Name())
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("[error] failed to execute command from memory filesystem:", err)
	}
	fmt.Println("cmd output:", output)

	stdout, err := self.MemFD.Command().Output()
	if err != nil {
		fmt.Println("[error] failed to execute command from memory fd:", err)
	}

	fmt.Println("cmd output:", string(stdout))

	//_, _, err1 = syscall.RawSyscall(syscall.SYS_EXECVE, 0, 0, 0)

	for {
	} // Hold open for testing

	return "", nil
}

func (self *Binary) PrintString() {
	fmt.Println("===[LOADED MemFD Executable ]================")
	fmt.Println("| |-[Executable Path] ", self.ExecutablePath)
	fmt.Println("| |-[Path] ", self.Path)
	fmt.Println("| |-[Name] ", self.Name)
	fmt.Println("| |-[Size] ", self.Size)
	fmt.Println("=============================================")
}
