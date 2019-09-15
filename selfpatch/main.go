package main

import (
	"fmt"

	"github.com/multiverse-os/singularity/test/binarypatch"
)

func main() {
	fmt.Println("Hello")

	myself := binarypatch.ReadMyself(500, -1, "linux")
	fmt.Println("Reading from myself: ", myself, string(myself))
}
