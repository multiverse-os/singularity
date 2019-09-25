package singularity

import (
	"fmt"
	"strings"

	//memexec "github.com/multiverse-os/singularity/memexec"
	memfd "github.com/multiverse-os/singularity/memfd"
)

type Binary struct {
	Size     int
	Output   string
	ExitCode int
	FD       *memfd.MemFD
}

func LoadExecutable(name string, bytes []byte) *Binary {
	binary := &Binary{
		FD: memfd.New(name),
	}
	bytesWritten, err := binary.FD.Write(bytes)
	if err != nil {
		fmt.Println("[error] failed to write data to fd:", err)
	}
	binary.Size = bytesWritten
	binary.String()
	return binary
}

func (self *Binary) Run(arguments ...string) error {
	fmt.Println("[singularity] ? Inside Execute():", self.FD.Name())

	self.FD.Execute(arguments...)
	//fmt.Println("pid:", pid)
	//fmt.Println("fd path:", fmt.Sprintf("/proc/self/fd/%d", fd))

	//outBytes, err := memexec.MemFD(self.FD.File, arguments...).CombinedOutput()
	//if err != nil {
	//	fmt.Println("[error] failed to execute memexec.MemFD() from singularity:", err)
	//}
	//fmt.Println("output:", string(outBytes))

	return nil
}

func (self *Binary) String() {
	fmt.Println("  ++================+========================++")
	fmt.Println("  ||   Attribute    |          Value         ||")
	fmt.Println("  ++================+========================++")
	fmt.Println("  |      Name       |    ", self.FD.Name(), strings.Repeat(" ", (18-len(self.FD.Name()))), "|")
	fmt.Println("  +-----------------+-------------------------+")
	fmt.Println("  |      Size       |    ", self.Size, strings.Repeat(" ", (16-len(string(self.Size)))), "|")
	fmt.Println("  +-----------------+-------------------------+")
	fmt.Println("  |      Path       |    ", self.FD.Path(), strings.Repeat(" ", (18-len(self.FD.Path()))), "|")
	fmt.Println("  +-----------------+-------------------------+\n")
}
