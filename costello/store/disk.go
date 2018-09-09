package store

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr/costello/parser"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

var _ Store = &Disk{}

type Disk struct {
	DBPath    string
	DBPackage string
	global    map[string]*parser.File
}

func NewDisk(path string, pkg string) *Disk {
	if len(path) == 0 {
		path = filepath.Join("internal", "packr-packed")
	}
	if len(pkg) == 0 {
		pkg = "packed"
	}
	return &Disk{
		DBPath:    path,
		DBPackage: pkg,
		global:    map[string]*parser.File{},
	}
}

func (d *Disk) FileNames(box *parser.Box) ([]string, error) {
	path := box.AbsPath
	if len(box.AbsPath) == 0 {
		path = box.Path
	}
	var names []string
	err := godirwalk.Walk(path, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback: func(path string, de *godirwalk.Dirent) error {
			if !de.IsRegular() {
				return nil
			}
			names = append(names, path)
			return nil
		},
	})
	return names, err
}

func (d *Disk) Files(box *parser.Box) ([]*parser.File, error) {
	var files []*parser.File
	names, err := d.FileNames(box)
	if err != nil {
		return files, errors.WithStack(err)
	}
	for _, n := range names {
		b, err := ioutil.ReadFile(n)
		if err != nil {
			return files, errors.WithStack(err)
		}
		f := parser.NewFile(n, bytes.NewReader(b))
		files = append(files, f)
	}
	return files, nil
}

func (d *Disk) Pack(box *parser.Box) error {
	panic("not implemented")
}

func (d *Disk) Clean(box *parser.Box) error {
	root := box.PackageDir
	if len(root) == 0 {
		return errors.New("can't clean an empty box.PackageDir")
	}
	return Clean(root)
}

func (d *Disk) Close() error {
	fmt.Println("not implemented")
	return nil
}

func Clean(root string) error {
	if len(root) == 0 {
		pwd, err := os.Getwd()
		if err != nil {
			return errors.WithStack(err)
		}
		root = pwd
	}
	callback := func(path string, info *godirwalk.Dirent) error {
		base := filepath.Base(path)
		if base == ".git" || base == "vendor" || base == "node_modules" {
			return filepath.SkipDir
		}
		if info == nil || info.IsDir() {
			return nil
		}
		if strings.Contains(base, "-packr.go") {
			err := os.Remove(path)
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	}
	return godirwalk.Walk(root, &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	})
}

// resolve file paths (only) for the boxes
// compile "global" db
// resolve files for boxes to point at global db
// write global db to disk (default internal/packr)
// write boxes db to disk (default internal/packr)
// write -packr.go files in each package (1 per package) that init the global db
