package store

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"strings"

	"github.com/gobuffalo/packr/encoding/hex"

	"github.com/pkg/errors"
)

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
