package packr

import (
	"bytes"
	"net/http"
	"time"
)

// HTTPBox implements http.FileSystem which allows the use of Box with a http.FileServer.
//   e.g.: http.Handle("/", http.FileServer(packr.NewBox("http-files").HTTPBox()))
type HTTPBox struct {
	Box
}

// HTTPBox creates a new HTTPBox from an existing Box
func (b Box) HTTPBox() HTTPBox {
	return HTTPBox{
		Box: b,
	}
}

// Open returns a File using the http.File interface
func (hb HTTPBox) Open(name string) (http.File, error) {
	bb := &bytes.Buffer{}
	b, err := hb.MustBytes(name)
	if err != nil {
		return nil, err
	}
	bb.Write(b)
	f := file{
		Buffer: bb,
		Name:   name,
		info: fileInfo{
			Path:     name,
			Contents: b,
			size:     int64(len(b)),
			modTime:  time.Now(),
		},
	}
	return f, nil
}
