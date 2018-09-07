package packr

import (
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
	b.callingDir = resolver.Ident(cd)
	return b
}

// PackBytes packs bytes for a file into a box.
func PackBytes(box string, name string, bb []byte) {
	panic("not implemented")
	// gil.Lock()
	// defer gil.Unlock()
	// if _, ok := data[box]; !ok {
	// 	data[box] = map[string][]byte{}
	// }
	// data[box][name] = bb
}

// PackBytesGzip packets the gzipped compressed bytes into a box.
func PackBytesGzip(box string, name string, bb []byte) error {
	panic("not implemented")
	// var buf bytes.Buffer
	// w := gzip.NewWriter(&buf)
	// _, err := w.Write(bb)
	// if err != nil {
	// 	return err
	// }
	// err = w.Close()
	// if err != nil {
	// 	return err
	// }
	// PackBytes(box, name, buf.Bytes())
	// return nil
}

// PackJSONBytes packs JSON encoded bytes for a file into a box.
func PackJSONBytes(box string, name string, jbb string) error {
	panic("not implemented")
	// var bb []byte
	// err := json.Unmarshal([]byte(jbb), &bb)
	// if err != nil {
	// 	return err
	// }
	// PackBytes(box, name, bb)
	// return nil
}
