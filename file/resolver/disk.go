package resolver

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr/file"
)

var _ Resolver = &Disk{}

type Disk struct {
	Root Ident
}

func (d *Disk) Find(name Ident) (file.File, error) {
	path := filepath.Join(d.Root.OsPath(), name.OsPath())
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return file.NewDir(name.OsPath()), nil
	}
	if bb, err := ioutil.ReadFile(path); err == nil {
		return file.NewFile(name.Name(), bb), nil
	}
	return nil, os.ErrNotExist
}
