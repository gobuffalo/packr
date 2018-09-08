package packr

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
)

// File has been deprecated and file.File should be used instead
type File = file.File

var (
	// ErrResOutsideBox gets returned in case of the requested resources being outside the box
	ErrResOutsideBox = errors.New("can't find a resource outside the box")
)

// NewBox returns a Box that can be used to
// retrieve files from either disk or the embedded
// binary.
func NewBox(path string) Box {
	b := New(path)
	var cd string
	if !filepath.IsAbs(path) {
		_, filename, _, _ := runtime.Caller(1)
		cd = filepath.Dir(filename)
	}

	// this little hack courtesy of the `-cover` flag!!
	cov := filepath.Join("_test", "_obj_test")
	cd = strings.Replace(cd, string(filepath.Separator)+cov, "", 1)
	if !filepath.IsAbs(cd) && cd != "" {
		cd = filepath.Join(GoPath(), "src", cd)
	}
	b.ResolutionDir = resolver.Ident(cd)
	return b
}

// PackBytes packs bytes for a file into a box.
func PackBytes(box string, name string, bb []byte) {
	d := resolver.NewInMemory(map[resolver.Ident]file.File{})
	iname := resolver.Ident(name)
	d.Pack(iname, file.NewFile(iname.OsPath(), bb))
	resolver.Register(resolver.Ident(box), iname, d)
}

// PackBytesGzip packets the gzipped compressed bytes into a box.
func PackBytesGzip(box string, name string, bb []byte) error {
	// TODO: this function never did what it was supposed to do!
	PackBytes(box, name, bb)
	return nil
}

// PackJSONBytes packs JSON encoded bytes for a file into a box.
func PackJSONBytes(box string, name string, jbb string) error {
	var bb []byte
	err := json.Unmarshal([]byte(jbb), &bb)
	if err != nil {
		return err
	}
	PackBytes(box, name, bb)
	return nil
}
