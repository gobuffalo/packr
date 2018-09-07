package box

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"io"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	box := packr.NewBox("../box/templates")
	code, err := box.MustString("generated.go.plush")
	if err != nil {
		return g, errors.WithStack(err)
	}

	f := genny.NewFile(filepath.Join(opts.Root, opts.Name+"-packr.go.plush"), strings.NewReader(code))

	ctx := plush.NewContext()
	ctx.Set("opts", opts)

	files := []packedFiles{}

	for _, f := range opts.Files {
		packed, err := hexGzip(f.String())
		if err != nil {
			return g, errors.WithStack(err)
		}
		files = append(files, packedFiles{
			Box:    opts.Name,
			Name:   f.Name(),
			Packed: packed,
		})
	}

	ctx.Set("files", files)

	g.Transformer(plushgen.Transformer(ctx))
	g.Transformer(gotools.FmtTransformer())

	g.File(f)
	return g, nil
}

type packedFiles struct {
	Box    string
	Name   string
	Packed string
}

func hexGzip(s string) (string, error) {
	bb := &bytes.Buffer{}
	enc := hex.NewEncoder(bb)
	zw := gzip.NewWriter(enc)
	io.Copy(zw, strings.NewReader(s))
	zw.Close()

	return bb.String(), nil
}
