package packr

import (
	"testing"

	"github.com/gobuffalo/packr/file"
	"github.com/gobuffalo/packr/file/resolver"
	"github.com/stretchr/testify/require"
)

func Test_Box_AddString(t *testing.T) {
	r := require.New(t)

	resolver.ClearRegistry()

	box := New("./templates")
	s, err := box.MustString("foo.txt")
	r.Error(err)
	r.Equal("", s)

	r.NoError(box.AddString("foo.txt", "foo!!"))
	s, err = box.MustString("foo.txt")
	r.NoError(err)
	r.Equal("foo!!", s)
}

func Test_Box_String(t *testing.T) {
	r := require.New(t)

	resolver.ClearRegistry()
	d := resolver.NewInMemory(map[resolver.Ident]file.File{
		"foo.txt": file.NewFile("foo.txt", []byte("foo!")),
	})
	resolver.Register("./templates", "foo.txt", d)

	box := New("./templates")

	s := box.String("foo.txt")
	r.Equal("foo!", s)

	s = box.String("idontexist")
	r.Equal("", s)
}

func Test_Box_String_Miss(t *testing.T) {
	r := require.New(t)

	resolver.ClearRegistry()

	box := New("./_fixtures/templates")

	s := box.String("foo.txt")
	r.Equal("FOO!!!\n", s)

	s = box.String("idontexist")
	r.Equal("", s)
}
