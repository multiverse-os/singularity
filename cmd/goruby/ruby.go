package main

import (
	"fmt"
	"io/ioutil"

	memexec "github.com/multiverse-os/portalgun/go-memexec"
)

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
	rubyExecutable, err := memexec.Asset("/bin/ruby")
	if err != nil {
		return
	}
	//p, err := exec.LookPath("echo")
	//if err != nil {
	//	t.Fatal(err)
	//}

	binary, err := ioutil.ReadFile(rubyExecutable)
	if err != nil {
		t.Fatal(err)
	}

	m, err := memexec.New(binary)
	if err != nil {
		return
	}

	defer func() {
		e := m.Close()
		if err == nil && e != nil {
			t.Fatal(e)
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
