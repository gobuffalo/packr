package packr

import (
	"encoding/json"
	"errors"

	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/file/resolver"
)

// File has been deprecated and file.File should be used instead
type File = file.File

var (
	// ErrResOutsideBox gets returned in case of the requested resources being outside the box
	// Deprecated
	ErrResOutsideBox = errors.New("can't find a resource outside the box")
)

// PackBytes packs bytes for a file into a box.
// Deprecated
func PackBytes(box string, name string, bb []byte) {
	b := NewBox(box)
	d := resolver.NewInMemory(map[string]file.File{})
	f, err := file.NewFile(name, bb)
	if err != nil {
		panic(err)
	}
	if err := d.Pack(name, f); err != nil {
		panic(err)
	}
	b.SetResolver(name, d)
}

// PackBytesGzip packets the gzipped compressed bytes into a box.
// Deprecated
func PackBytesGzip(box string, name string, bb []byte) error {
	// TODO: this function never did what it was supposed to do!
	PackBytes(box, name, bb)
	return nil
}

// PackJSONBytes packs JSON encoded bytes for a file into a box.
// Deprecated
func PackJSONBytes(box string, name string, jbb string) error {
	var bb []byte
	err := json.Unmarshal([]byte(jbb), &bb)
	if err != nil {
		return err
	}
	PackBytes(box, name, bb)
	return nil
}
