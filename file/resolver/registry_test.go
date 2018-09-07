package resolver

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gobuffalo/packr/file"
	"github.com/stretchr/testify/require"
)

func Test_Registery(t *testing.T) {
	ClearRegistry()
	r := require.New(t)

	d := NewInMemory(map[Ident]file.File{
		"foo.txt": file.NewFile("foo.txt", []byte("foo!")),
	})

	Register("mybox", "foo.txt", d)

	f, err := Resolve("mybox", "foo.txt")
	r.NoError(err)
	fi, err := f.FileInfo()
	r.NoError(err)
	r.Equal("foo.txt", fi.Name())
	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("foo!", string(b))
}

func Test_Registery_Miss(t *testing.T) {
	ClearRegistry()
	r := require.New(t)

	f, err := Resolve(Ident("_fixtures/templates"), "foo.txt")
	r.NoError(err)
	fi, err := f.FileInfo()
	r.NoError(err)
	r.Equal("foo.txt", fi.Name())
	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("foo!", strings.TrimSpace(string(b)))
}
