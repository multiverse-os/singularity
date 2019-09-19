package main

import (
	"fmt"

	patch "github.com/multiverse-os/singularity/patch"
)

func main() {
	fmt.Println("binarypatch")
	fmt.Println("=======================")
	selfPatch, err := patch.ReadSelf(5, -1)
	if err != nil {
		fmt.Println("[error] self patch failed to parse elf:", selfPatch)
	}
	fmt.Println("selfPatch:", selfPatch)

	//b := binarypatch.New("../sample/sample_darwin64", "linux")
	//index := b.Locate()
	//b.Write([]byte("hello"), index)

	//b.WriteFile("patched")

	//b2 := binarypatch.New("./patched", "darwin")
	//index2 := b2.Locate()
	//data := b2.Read(index2, 5)

	//fmt.Println("Stuff: ", data, string(data))

	//myself := binarypatch.ReadMyself(5, -1, "darwin")

	//fmt.Println("Reading from myself: ", myself)

	//b = binarypatch.New("../sample/sample_win64.exe", "windows")
	//index = b.Locate()
	//b.Write([]byte("hello"), index)
	//b.WriteFile("patched.exe")

}
