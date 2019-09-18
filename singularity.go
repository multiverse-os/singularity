package singularity

import (
	"fmt"
	"strings"

	memfd "github.com/multiverse-os/singularity/memfd"
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

const (
	mfdCloexec  = 0x0001
	memfdCreate = 319
)

type Binary struct {
	Name     string
	Size     int
	Output   string
	ExitCode int
	FD       *memfd.MemFD
}

func NewBinary(name string, bytes []byte) *Binary {
	binary := &Binary{
		Name: name,
		FD:   memfd.New(name),
	}
	bytesWritten, err := binary.FD.Write(bytes)
	if err != nil {
		fmt.Println("[error] failed to write data to fd:", err)
	}
	binary.Size = bytesWritten
	binary.String()
	return binary
}

func (self *Binary) Execute(arguments string) error {
	fmt.Println("[singularity] Inside Execute():", self.Name)
	return self.FD.Exec(arguments)
}

func (self *Binary) String() {
	fmt.Println("  ++================+========================++")
	fmt.Println("  ||   Attribute    |          Value         ||")
	fmt.Println("  ++================+========================++")
	fmt.Println("  |      Name       |    ", self.Name, strings.Repeat(" ", (18-len(self.Name))), "|")
	fmt.Println("  +-----------------+-------------------------+")
	fmt.Println("  |      Size       |    ", self.Size, strings.Repeat(" ", (16-len(string(self.Size)))), "|")
	fmt.Println("  +-----------------+-------------------------+")
	fmt.Println("  |      Path       |    ", self.FD.Path(), strings.Repeat(" ", (18-len(self.FD.Path()))), "|")
	fmt.Println("  +-----------------+-------------------------+\n")
}
