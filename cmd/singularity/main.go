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

	binary := singularity.NewBinary("ruby", binaries.Ruby)

	binary.Execute("-e p 'test ruby code'")

}
