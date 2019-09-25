package main

import (
	"fmt"
	//"os"

	singularity "github.com/multiverse-os/singularity"
	//memexec "github.com/multiverse-os/singularity/memexec"
	executable "github.com/multiverse-os/singularity/store/executable"
	//vfs "github.com/multiverse-os/singularity/store/vfs"
	//memfs "github.com/multiverse-os/singularity/store/vfs/memfs"
	//mountfs "github.com/multiverse-os/singularity/store/vfs/mountfs"
)

func main() {
	fmt.Println("[singularity]: memory execution of binary")
	fmt.Println("============================================================")
	fmt.Println("An example of binary execution completely in memory without ")
	fmt.Println("touching the disk, or creating temporary files using memFD.\n")

	err := singularity.LoadExecutable("ruby", executable.Ruby).Run("-e 'p \"hello world from ruby\"'")
	if err != nil {
		fmt.Println("[error] failed to run executable from memory:", err)
	}

	//osfs := vfs.OS()
	//osfs.Mkdir("/tmp", 0777)
	//fs := mountfs.Create(osfs)
	//mfs := memfs.Create()
	//mfs.Mkdir("/memfs", 0777)
	//fs.Mount(mfs, "/memfs")
	//fs.Mkdir("/memfs/rubyscripts", 0777)

	//memoryFile, _ := osfs.OpenFile("/memfs/rubyscripts/helloworld.rb", os.O_RDWR, 0)
	//_, err := memoryFile.Write([]byte("p 'hello world from ruby temp file'"))
	//if err != nil {
	//	fmt.Errorf("[error] failed to write to memory filesystem:\n", err)
	//}

	//memoryCommand.Run("-e 'hello world from ruby'")

	//outputBytes, err := memexec.Command(executable.Ruby, " /memfs/rubyscripts/helloworld.rb").Run()
	//if err != nil {
	//	fmt.Println("[error] failed to execute binary from memory:", err)
	//}
	//fmt.Println("output bytes:", string(outputBytes))

	//err = memoryCommand.Run("-e p 'test ruby code'")
	//if err != nil {
	//	fmt.Println("[error] failed to execute binary from memory:", err)
	//}

}
