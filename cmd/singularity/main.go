package main

import (
	"fmt"

	singularity "github.com/multiverse-os/singularity"
	binaries "github.com/multiverse-os/singularity/binaries"
)

func main() {
	fmt.Println("[singularity]: memory execution of binary")
	fmt.Println("============================================================")
	fmt.Println("An example of binary execution completely in memory without ")
	fmt.Println("touching the disk, or creating temporary files using memFD.\n")

	binary := singularity.LoadBinary("ruby", binaries.Ruby)

	output, err := binary.Execute("test.rb")
	if err != nil {
		fmt.Println("[singularity] failed to execute binary in memory:", err)
	}

	fmt.Println("output:", string(output))
}
