package binaries

import (
	"os"
)

type StoredData interface {
	Encode() []byte
	Decode() []byte
	Disk() *os.File
	Memory() *os.File
}

type CompressionAlgorithm int

const (
	Zstd CompressionAlgorithm = iota
	Snappy
	//Brotli
)

func Compress() []byte {
	return []byte{}
}

func Extract() []byte {
	return []byte{}
}
