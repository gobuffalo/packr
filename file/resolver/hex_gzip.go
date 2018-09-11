package resolver

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/gobuffalo/packr/file"
	"github.com/pkg/errors"
)

var _ Resolver = &HexGzip{}

type HexGzip struct {
	packed   map[string]string
	unpacked map[string]file.File
	moot     *sync.RWMutex
}

var _ file.FileMappable = &HexGzip{}

func (hg *HexGzip) FileMap() map[string]file.File {
	hg.moot.RLock()
	var names []string
	for k := range hg.packed {
		names = append(names, k)
	}
	hg.moot.RUnlock()
	m := map[string]file.File{}
	for _, n := range names {
		if f, err := hg.Find("", n); err == nil {
			m[n] = f
		}
	}
	return m
}

func (hg *HexGzip) Find(box string, name string) (file.File, error) {
	fmt.Println("HexGzip: Find", name)
	hg.moot.RLock()
	if f, ok := hg.unpacked[name]; ok {
		hg.moot.RUnlock()
		return f, nil
	}
	hg.moot.RUnlock()
	packed, ok := hg.packed[name]
	if !ok {
		return nil, os.ErrNotExist
	}

	unpacked, err := UnGzip(packed)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	f := file.NewFile(OsPath(name), []byte(unpacked))
	hg.moot.Lock()
	hg.unpacked[name] = f
	hg.moot.Unlock()
	return f, nil
}

func NewHexGzip(files map[string]string) (*HexGzip, error) {
	if files == nil {
		files = map[string]string{}
	}
	hg := &HexGzip{
		packed:   files,
		unpacked: map[string]file.File{},
		moot:     &sync.RWMutex{},
	}

	return hg, nil
}

func Gzip(s string) (string, error) {
	bb := &bytes.Buffer{}
	var w io.Writer = bb
	w, err := gzip.NewWriterLevel(w, gzip.BestCompression)
	if err != nil {
		return "", errors.WithStack(err)
	}
	io.Copy(w, strings.NewReader(s))
	if cl, ok := w.(io.Closer); ok {
		cl.Close()
	}

	encoded := base64.StdEncoding.EncodeToString(bb.Bytes())
	return encoded, nil
}

func UnGzip(packed string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(packed)
	if err != nil {
		return "", errors.WithStack(err)
	}
	br := bytes.NewBuffer(decoded)
	r, err := gzip.NewReader(br)
	if err != nil {
		return "", errors.WithStack(err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.WithStack(err)
	}
	r.Close()
	return string(b), nil
}
