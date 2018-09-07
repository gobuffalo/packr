package box

import (
	"context"
	"strings"
	"testing"

	"github.com/gobuffalo/genny"
	"github.com/stretchr/testify/require"
)

func Test_New(t *testing.T) {
	r := require.New(t)

	g, err := New(&Options{
		Name:    "../foo/templates",
		Package: "baz",
		Files: []genny.File{
			genny.NewFile("example.txt", strings.NewReader("hi!!!")),
		},
		Root: "./foo",
	})
	r.NoError(err)

	run := genny.DryRunner(context.Background())
	run.With(g)

	r.NoError(run.Run())

	res := run.Results()

	r.Len(res.Commands, 0)
	r.Len(res.Files, 1)

	f := res.Files[0]
	r.Equal("foo/templates-packr.go", f.Name())
	body := f.String()
	r.Contains(body, `packr.PackHexGzip("../foo/templates", "example.txt", "1f8b08000000000000ffcac854545404040000ffff725a37d005000000")`)
	r.Contains(body, "package baz")
}
