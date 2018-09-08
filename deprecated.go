package packr

import (
	"encoding/json"
	"errors"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
)

// File has been deprecated and file.File should be used instead
type File = file.File

var (
	// ErrResOutsideBox gets returned in case of the requested resources being outside the box
	ErrResOutsideBox = errors.New("can't find a resource outside the box")
)

// PackBytes packs bytes for a file into a box.
func PackBytes(box string, name string, bb []byte) {
	b := NewBox(box)
	d := resolver.NewInMemory(map[resolver.Ident]file.File{})
	if err := d.Pack(resolver.Ident(name), file.NewFile(name, bb)); err != nil {
		panic(err)
	}
	b.SetResolver(name, d)
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
