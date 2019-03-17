package main

import (
	"fmt"
	"io/ioutil"
	"os"

	memexec "github.com/multiverse-os/singularity"
)

var rubyBinaryPath = "/usr/bin/ruby"

func main() {
	fmt.Println("Ruby embedded in Go")
	fmt.Println("===================")
	output, err := RubyExec("test.rb")
	if err != nil {
		fmt.Println("[error] failed to execute ruby binary:", err)
	} else {
		fmt.Println("Executing....", output)
	}
}

func RubyExec(argv ...string) ([]byte, error) {
	// Asset function provided by go-bindata
	//rubyExecutable, err := memexec.Asset("/bin/ruby")
	//if err != nil {
	//	return nil, err
	//}
	//p, err := exec.LookPath("echo")
	//if err != nil {
	//	t.Fatal(err)
	//}
	if fileInfo, err := os.Stat(rubyBinaryPath); os.IsNotExist(err) && !fileInfo.IsDir() {
		fmt.Println("[fatal error] binary '"+rubyBinaryPath+"' is missing:", err)
		// TODO: If missing, send to manual input
	} else {

		dataAsBytes, err := ioutil.ReadFile(rubyBinaryPath)
		if err != nil {
			fmt.Errorf("[error] failed to read [ Ruby ] binary:", err)
			return []byte{}, err
		} else {
			fmt.Println("[singularity] successfully loaded binary file into memory")
			fmt.Println("[singularity] loaded [", len(dataAsBytes), "b ] into memory")
			fmt.Println("[singularity] Name()?", fileInfo.Name())
			fmt.Println("[singularity] Size()?", fileInfo.Size())
			fmt.Println("[singularity] Mode()?", fileInfo.Mode())
			fmt.Println("[singularity] ModTime()?", fileInfo.ModTime())
			fmt.Println("[singularity] IsDir()?", fileInfo.IsDir())
			fmt.Println("[singularity] Sys()?", fileInfo.Sys())

			fmt.Println("\n[singularity] now attempting to execute the binary file...\n")
		}

		m, err := memexec.New(fileInfo.Sys().([]byte))
		if err != nil {
			return []byte{}, fmt.Errorf("[error] failed to execute binary ["+rubyBinaryPath+"]:", err)
		}

		defer func() {
			e := m.Close()
			if err == nil && e != nil {
				fmt.Errorf("[error] failed to close binary file:", err)
			}
		}()

		// m can be cached to avoid extra copying
		// when it's needed exec the same code multiple times
		//defer func() {
		//	cerr := m.Close()
		//	if err == nil {
		//		err = cerr
		//	}
		//}()

		return m.Command(argv...).Output()
	}
	return nil, nil
}
