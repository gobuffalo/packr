package file

import (
	"bytes"
	"io"
	"os"
)

// File represents a virtual, or physical, backing of
// a file object in a Box
type File interface {
	Name() string
	io.ReadCloser
	io.Writer
	FileInfo() (os.FileInfo, error)
	Readdir(count int) ([]os.FileInfo, error)
	Seek(offset int64, whence int) (int64, error)
	Stat() (os.FileInfo, error)
}

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
