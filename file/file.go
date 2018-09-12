package file

import (
	"bytes"
	"io"
	"os"
)

type File interface {
	Name() string
	io.ReadCloser
	io.Writer
	FileInfo() (os.FileInfo, error)
	Readdir(count int) ([]os.FileInfo, error)
	Seek(offset int64, whence int) (int64, error)
	Stat() (os.FileInfo, error)
}

type FileMappable interface {
	FileMap() map[string]File
}

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
