package packr

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
)

// NewBox returns a Box that can be used to
// retrieve files from either disk or the embedded
// binary.
func NewBox(path string) Box {
	_, filename, _, _ := runtime.Caller(1)
	return Box{
		Path:       path,
		callingDir: filepath.Dir(filename),
	}
}

// Box represent a folder on a disk you want to
// have access to in the built Go binary.
type Box struct {
	Path       string
	callingDir string
	data       map[string][]byte
}

// String of the file asked for or an empty string.
func (b Box) String(name string) string {
	return string(b.Bytes(name))
}

// MustString returns either the string of the requested
// file or an error if it can not be found.
func (b Box) MustString(name string) (string, error) {
	bb, err := b.MustBytes(name)
	return string(bb), err
}

// Bytes of the file asked for or an empty byte slice.
func (b Box) Bytes(name string) []byte {
	bb, _ := b.MustBytes(name)
	return bb
}

// MustBytes returns either the byte slice of the requested
// file or an error if it can not be found.
func (b Box) MustBytes(name string) ([]byte, error) {
	bb, err := find(b.Path, name)
	if err == nil {
		return bb, err
	}
	p := filepath.Join(b.callingDir, b.Path, name)
	return ioutil.ReadFile(p)
}
