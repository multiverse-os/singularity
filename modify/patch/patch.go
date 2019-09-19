package patch

import (
	"bytes"
	"io/ioutil"
	"os"
)

type Arch int

const (
	ELF64 Arch = iota
	ELF32
)

// TODO: Just detect the arch
type Patch struct {
	Arch     Arch
	Filename string
	Data     []byte
}

// TODO: Cant we just take filename from the current executalbe?
func New(filename string) (*Patch, error) {
	patch := &Patch{
		Filename: filename,
		Arch:     ELF64,
	}
	var err error
	patch.Data, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return patch, nil
}

func ReadSelf(length int, index int) ([]byte, error) {
	filename, err := os.Executable()
	if err != nil {
		return []byte{}, err
	}
	patch, err := New(filename)
	if index == -1 {
		index = patch.Locate()
	}
	return patch.Read(index, length), nil
}

func (self *Patch) Read(index int, length int) []byte { return self.Data[index : index+length] }
func (self *Patch) PatchFile() error                  { return ioutil.WriteFile(self.Filename, self.Data, 0744) }

func (self *Patch) Locate() (index int) {
	str := []byte("__TEXT")
	for i, v := range self.Data {
		n := 0
		byteGroup := []byte{v}
		for n <= len(str) && i < len(self.Data)-len(str) {
			byteGroup = append(byteGroup, self.Data[i+n])
			n++
		}
		if bytes.Index(byteGroup, str) == 0 {
			index = i + len(str)
			break
		}
	}
	return index
}

func (self *Patch) Write(data []byte, index int) {
	for i, v := range data {
		position := index + i
		self.Data[position] = v
	}
}
