package file

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

var virtualFileModTime = time.Now()
var _ File = virtualFile{}

type virtualFile struct {
	*bytes.Reader
	name string
	info info
}

func (f virtualFile) Name() string {
	return f.name
}

func (f virtualFile) FileInfo() (os.FileInfo, error) {
	return f.info, nil
}

func (f virtualFile) Close() error {
	return nil
}

func (f virtualFile) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("not implemented")
}

func (f virtualFile) Readdir(count int) ([]os.FileInfo, error) {
	return []os.FileInfo{f.info}, nil
}

func (f virtualFile) Stat() (os.FileInfo, error) {
	return f.info, nil
}
