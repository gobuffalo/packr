package packr

import (
	"errors"

	"github.com/gobuffalo/packr/file"
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
	// TODO: deprecate
	return New(path)
	// var cd string
	// if !filepath.IsAbs(path) {
	// 	_, filename, _, _ := runtime.Caller(1)
	// 	cd = filepath.Dir(filename)
	// }
	//
	// // this little hack courtesy of the `-cover` flag!!
	// cov := filepath.Join("_test", "_obj_test")
	// cd = strings.Replace(cd, string(filepath.Separator)+cov, "", 1)
	// if !filepath.IsAbs(cd) && cd != "" {
	// 	cd = filepath.Join(GoPath(), "src", cd)
	// }
	//
	// return Box{
	// 	Path:       path,
	// 	callingDir: cd,
	// 	data:       map[string][]byte{},
	// }
}
