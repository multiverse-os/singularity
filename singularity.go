package singularity

// NOTE: ALternative way to update,
// store the soruce code, use it to rebuild and replace itself

import (
	"fmt"
	"os"

	memfd "github.com/multiverse-os/singularity/memfd"
	memfs "github.com/multiverse-os/singularity/memfs"
	msyscall "github.com/multiverse-os/singularity/msyscall"
)

type Process struct {
	PID int
}

type Binary struct {
	ParentProcess Process
	Name          string
	Size          int
	Data          []byte
	Output        string
	ExitCode      int
	MemFD         memfd.MemFD
}

func LoadBinary(name string, data []byte) *Binary {
	binary := &Binary{
		ParentProcess: Process{
			PID: os.Getpid(),
		},
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
	fmt.Println("[singularity] fd path:", self.MemFD.FDPath())
	return self
}

func (self *Binary) Execute() (string, error) {
	fmt.Println("[singularity] Inside Execute()")

	//filepath := fmt.Sprintf("/proc/%d/fd/%d", pid, self.MemFD)
	//fmt.Println("filepath with pid:[", pid, "] :", filepath)

	//f, err := os.OpenFile(self.MemFD.FDPath(), os.O_RDWR, 0755)
	perm, err := self.MemFD.FDPerm()
	if err != nil {
		fmt.Println("[error] failed to obtain permissions for memory fd:", err)
	}
	fmt.Println("fd permissions are:", perm.String())

	//self.MemFD.SetCloexec()

	fs, _ := memfs.NewFS() // remember kids don't ignore errors

	// Opens a file with read/write permissions in the current directory
	f, _ := fs.Create("/example.txt")

	f.Write([]byte("Hello, world!"))
	f.Close()

	stdout, err := self.MemFD.Command().Output()
	if err != nil {
		fmt.Println("[error] failed to execute command from memory:", err)
	}

	fmt.Println("cmd output:", string(stdout))

	//_, _, err1 = syscall.RawSyscall(syscall.SYS_EXECVE, 0, 0, 0)

	for {
	} // Hold open for testing

	return "", nil
}
