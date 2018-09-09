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
	Root string
}

func (d *Disk) Find(box string, name string) (file.File, error) {
	fmt.Println("Disk: Find", name)
	path := OsPath(name)
	if !filepath.IsAbs(path) {
		path = filepath.Join(OsPath(d.Root), path)
	}
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if fi.IsDir() {
		return file.NewDir(OsPath(name)), nil
	}
	if bb, err := ioutil.ReadFile(path); err == nil {
		return file.NewFile(OsPath(name), bb), nil
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
		name := strings.TrimPrefix(path, OsPath(d.Root)+string(filepath.Separator))
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		m[name] = file.NewFile(name, b)
		moot.Unlock()
		return nil
	}
	err := godirwalk.Walk(OsPath(d.Root), &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	})
	if err != nil {
		fmt.Println("error walking", OsPath(d.Root), err)
	}
	return m
}
