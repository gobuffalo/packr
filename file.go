package packr

import (
	"bytes"
	"os"
)

type file struct {
	*bytes.Buffer
	Name string
	info fileInfo
}

func (f file) Close() error {
	return nil
}

func (f file) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func (f file) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{f.info}, nil
}

func (f file) Stat() (os.FileInfo, error) {
	return f.info, nil
}
