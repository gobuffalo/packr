package store

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gobuffalo/packr/costello/parser"
	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

var _ Store = &Disk{}

type Disk struct {
	DBPath    string
	DBPackage string
	global    map[string]string
	boxes     map[string]*parser.Box
}

func NewDisk(path string, pkg string) *Disk {
	if len(path) == 0 {
		path = "packr-packed"
	}
	if len(pkg) == 0 {
		pkg = "packed"
	}
	return &Disk{
		DBPath:    path,
		DBPackage: pkg,
		global:    map[string]string{},
		boxes:     map[string]*parser.Box{},
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
	d.boxes[box.Name] = box
	names, err := d.FileNames(box)
	if err != nil {
		return errors.WithStack(err)
	}
	for _, n := range names {
		k, ok := d.global[n]
		if ok {
			continue
		}
		k = makeKey(box, n)
		// not in the global, so add it!
		d.global[n] = k
	}
	return nil
}

func (d *Disk) Clean(box *parser.Box) error {
	root := box.PackageDir
	if len(root) == 0 {
		return errors.New("can't clean an empty box.PackageDir")
	}
	return Clean(root)
}

type options struct {
	Package     string
	GlobalFiles map[string]string
	Boxes       []optsBox
}

type optsBox struct {
	Name string
	Path string
}

func (d *Disk) Close() error {
	opts := options{
		Package:     d.DBPackage,
		GlobalFiles: map[string]string{},
	}
	for k, v := range d.global {
		err := func() error {
			bb := &bytes.Buffer{}
			enc := hex.NewEncoder(bb)
			zw := gzip.NewWriter(enc)
			f, err := os.Open(k)
			if err != nil {
				return errors.WithStack(err)
			}
			defer f.Close()
			io.Copy(zw, f)
			if err := zw.Close(); err != nil {
				return errors.WithStack(err)
			}
			opts.GlobalFiles[v] = bb.String()
			return nil
		}()
		if err != nil {
			return errors.WithStack(err)
		}
	}
	for _, b := range d.boxes {
		ob := optsBox{
			Name: b.Name,
		}
		opts.Boxes = append(opts.Boxes, ob)
	}
	sort.Slice(opts.Boxes, func(a, b int) bool {
		return opts.Boxes[a].Name < opts.Boxes[b].Name
	})
	t, err := template.New("").Parse(diskGlobalTmpl)
	if err != nil {
		return errors.WithStack(err)
	}
	bb := &bytes.Buffer{}
	if err := t.Execute(bb, opts); err != nil {
		return errors.WithStack(err)
	}
	os.MkdirAll(d.DBPath, 0755)
	fp := filepath.Join(d.DBPath, "packed-packr.go")
	f, err := os.Create(fp)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	io.Copy(f, bb)
	// fmt.Println("### bb.String() ->", bb.String())
	return nil
}

// resolve file paths (only) for the boxes
// compile "global" db
// resolve files for boxes to point at global db
// write global db to disk (default internal/packr)
// write boxes db to disk (default internal/packr)
// write -packr.go files in each package (1 per package) that init the global db

func makeKey(box *parser.Box, path string) string {
	s := strings.TrimPrefix(path, box.AbsPath)
	return strings.TrimPrefix(s, string(filepath.Separator))
}
