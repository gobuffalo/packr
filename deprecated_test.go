package packr

import (
	"testing"

	"github.com/gobuffalo/packr/file/resolver"
	"github.com/stretchr/testify/require"
)

func Test_NewBox(t *testing.T) {
	resolver.ClearRegistry()

	r := require.New(t)
	b := NewBox("./_fixtures/list_test")
	r.Len(b.List(), 4)
}

func Test_PackBytes(t *testing.T) {
	resolver.ClearRegistry()

	r := require.New(t)

	box := New("my/box")
	name := "foo.txt"
	body := []byte("foo!!")
	PackBytes(box.Name.Name(), name, body)

	f, err := box.MustString(name)
	r.NoError(err)
	r.Equal(string(body), f)
}

func Test_PackBytesGzip(t *testing.T) {
	resolver.ClearRegistry()

	r := require.New(t)

	box := New("my/box")
	name := "foo.txt"
	body := []byte("foo!!")
	PackBytesGzip(box.Name.Name(), name, body)

	f, err := box.MustString(name)
	r.NoError(err)
	r.Equal(string(body), f)
}

func Test_PackJSONBytes(t *testing.T) {
	resolver.ClearRegistry()

	r := require.New(t)

	box := New("my/box")
	name := "foo.txt"
	body := "\"PGgxPnRlbXBsYXRlcy9tYWlsZXJzL2xheW91dC5odG1sPC9oMT4KCjwlPSB5aWVsZCAlPgo=\""
	PackJSONBytes(box.Name.Name(), name, body)

	f, err := box.MustString(name)
	r.NoError(err)
	r.Equal("<h1>templates/mailers/layout.html</h1>\n\n<%= yield %>\n", f)
}
