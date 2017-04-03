package packr

import (
	"bytes"
	"os"
	"time"
)

var virtualFileModTime = time.Now()
var _ File = virtualFile{}

type virtualFile struct {
	*bytes.Buffer
	Name string
	info fileInfo
}

func (v virtualFile) FileInfo() (os.FileInfo, error) {
	return v.info, nil
}

func (f virtualFile) Close() error {
	return nil
}

func (f virtualFile) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f virtualFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{f.info}, nil
}

func (f virtualFile) Stat() (os.FileInfo, error) {
	return f.info, nil
}

func newVirtualFile(name string, b []byte) File {
	return virtualFile{
		Buffer: bytes.NewBuffer(b),
		Name:   name,
		info: fileInfo{
			Path:     name,
			Contents: b,
			size:     int64(len(b)),
			modTime:  virtualFileModTime,
		},
	}
}
