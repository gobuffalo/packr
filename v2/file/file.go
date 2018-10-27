package file

import (
	"bytes"

	"github.com/gobuffalo/packd"
)

// File represents a virtual, or physical, backing of
// a file object in a Box
type File = packd.File

// FileMappable types are capable of returning a map of
// path => File
type FileMappable interface {
	FileMap() map[string]File
}

// NewFile returns a virtual File implementation
func NewFile(name string, b []byte) File {
	return virtualFile{
		Reader: bytes.NewReader(b),
		name:   name,
		info: info{
			Path:     name,
			Contents: b,
			size:     int64(len(b)),
			modTime:  virtualFileModTime,
		},
	}
}

// NewDir returns a virtual dir implementation
func NewDir(name string) File {
	var b []byte
	v := NewFile(name, b).(virtualFile)
	v.info.isDir = true
	return v
}
