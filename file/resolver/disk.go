package resolver

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gobuffalo/packr/file"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

var _ Resolver = &Disk{}

type Disk struct {
	Root Ident
}

func (d *Disk) Find(name Ident) (file.File, error) {
	fmt.Println("Disk: Find", name)
	path := name.OsPath()
	if !filepath.IsAbs(path) {
		path = filepath.Join(d.Root.OsPath(), path)
	}
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

var _ file.FileMappable = &Disk{}

func (d *Disk) FileMap() map[string]file.File {
	moot := &sync.Mutex{}
	m := map[string]file.File{}
	callback := func(path string, de *godirwalk.Dirent) error {
		if !de.IsRegular() {
			return nil
		}
		moot.Lock()
		name := strings.TrimPrefix(path, d.Root.OsPath()+string(filepath.Separator))
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		m[name] = file.NewFile(name, b)
		moot.Unlock()
		return nil
	}
	err := godirwalk.Walk(d.Root.OsPath(), &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	})
	if err != nil {
		fmt.Println("error walking", d.Root.OsPath(), err)
	}
	return m
}
