package resolver

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
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
	packed   map[Ident]string
	unpacked map[Ident]file.File
	moot     *sync.RWMutex
}

func (hg *HexGzip) Find(name Ident) (file.File, error) {
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

	unpacked, err := unHexGzip(packed)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	f := file.NewFile(name.Name(), []byte(unpacked))
	hg.moot.Lock()
	hg.unpacked[name] = f
	hg.moot.Unlock()
	return f, nil
}

func NewHexGzip(files map[Ident]string) (*HexGzip, error) {
	if files == nil {
		files = map[Ident]string{}
	}
	hg := &HexGzip{
		packed:   files,
		unpacked: map[Ident]file.File{},
		moot:     &sync.RWMutex{},
	}

	return hg, nil
}

func hexGzip(s string) (string, error) {
	bb := &bytes.Buffer{}
	enc := hex.NewEncoder(bb)
	zw := gzip.NewWriter(enc)
	io.Copy(zw, strings.NewReader(s))
	zw.Close()

	return bb.String(), nil
}

func unHexGzip(packed string) (string, error) {
	br := bytes.NewBufferString(packed)
	dec := hex.NewDecoder(br)
	zr, err := gzip.NewReader(dec)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer zr.Close()

	b, err := ioutil.ReadAll(zr)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return string(b), nil
}
