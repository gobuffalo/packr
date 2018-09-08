package resolver

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_HexGzip_Find(t *testing.T) {
	r := require.New(t)

	x, err := hexGzip("foo!")
	r.NoError(err)
	files := map[string]string{
		"foo.txt": x,
	}
	d, err := NewHexGzip(files)

	f, err := d.Find("foo.txt")
	r.NoError(err)

	fi, err := f.FileInfo()
	r.NoError(err)
	r.Equal("foo.txt", fi.Name())

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("foo!", strings.TrimSpace(string(b)))
}
