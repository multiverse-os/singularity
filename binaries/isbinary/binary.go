package isbinary

import (
	"path"
	"strings"

	"github.com/shamsher31/gobinext"
)

var extensions = binext.Get()

// Extensions is the extensions for different binary types
var Extensions map[string]bool

func init() {
	Extensions = map[string]bool{}
	for _, ext := range extensions {
		Extensions[strings.ToLower(ext)] = true
	}
}

// Is checks if a path is a binary
func Is(p string) bool {
	ext := strings.TrimLeft(path.Ext(p), ".")
	return Extensions[strings.ToLower(ext)]
}
