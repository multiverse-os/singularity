package main

import (
	"fmt"

	singularity "github.com/multiverse-os/singularity"
)

func main() {
	fmt.Println("memory execution of binary")
	fmt.Println("=============================================")
	binary := singularity.LoadBinary("helloworld", singularity.HelloworldBytes)
	fmt.Println("[singularity] attempting to execute binary")
	output, err := binary.Execute()
	if err != nil {
		fmt.Println("[error] failed to execute embedded binary from memory:", err)
	}
	fmt.Println("embedded binary output:", output)

	fmt.Println("exciting")
}
