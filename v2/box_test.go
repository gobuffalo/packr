package packr

import (
	"path/filepath"
	"testing"

	"github.com/gobuffalo/packr/v2/file"
	"github.com/gobuffalo/packr/v2/file/resolver"
	"github.com/stretchr/testify/require"
)

func Test_NewBox(t *testing.T) {
	r := require.New(t)

	box := NewBox(filepath.Join("_fixtures", "list_test"))
	r.Len(box.List(), 4)

}
func Test_Box_AddString(t *testing.T) {
	r := require.New(t)

	box := NewBox("./templates")
	s, err := box.MustString("foo.txt")
	r.Error(err)
	r.Equal("", s)

	r.NoError(box.AddString("foo.txt", "foo!!"))
	s, err = box.MustString("foo.txt")
	r.NoError(err)
	r.Equal("foo!!", s)
}

func Test_Box_AddBytes(t *testing.T) {
	r := require.New(t)

	box := NewBox("Test_Box_AddBytes")
	s, err := box.MustString("foo.txt")
	r.Error(err)
	r.Equal("", s)

	r.NoError(box.AddBytes("foo.txt", []byte("foo!!")))
	s, err = box.MustString("foo.txt")
	r.NoError(err)
	r.Equal("foo!!", s)
}

func Test_Box_String(t *testing.T) {
	r := require.New(t)

	box := NewBox("./templates")
	d := resolver.NewInMemory(map[string]file.File{
		"foo.txt": file.NewFile("foo.txt", []byte("foo!")),
	})
	box.SetResolver("foo.txt", d)

	s := box.String("foo.txt")
	r.Equal("foo!", s)

	s = box.String("idontexist")
	r.Equal("", s)
}

func Test_Box_String_Miss(t *testing.T) {
	r := require.New(t)

	box := NewBox(filepath.Join("_fixtures", "templates"))

	s := box.String("foo.txt")
	r.Equal("FOO!!!\n", s)

	s = box.String("idontexist")
	r.Equal("", s)
}

func Test_Box_MustString(t *testing.T) {
	r := require.New(t)

	box := NewBox("./templates")
	d := resolver.NewInMemory(map[string]file.File{
		"foo.txt": file.NewFile("foo.txt", []byte("foo!")),
	})
	box.SetResolver("foo.txt", d)

	s, err := box.MustString("foo.txt")
	r.NoError(err)
	r.Equal("foo!", s)

	s, err = box.MustString("idontexist")
	r.Error(err)
	r.Equal("", s)
}

func Test_Box_MustString_Miss(t *testing.T) {
	r := require.New(t)

	box := NewBox(filepath.Join("_fixtures", "templates"))

	s, err := box.MustString("foo.txt")
	r.NoError(err)
	r.Equal("FOO!!!\n", s)

	s, err = box.MustString("idontexist")
	r.Error(err)
	r.Equal("", s)
}

func Test_Box_Bytes(t *testing.T) {
	r := require.New(t)

	box := NewBox("./templates")
	d := resolver.NewInMemory(map[string]file.File{
		"foo.txt": file.NewFile("foo.txt", []byte("foo!")),
	})
	box.SetResolver("foo.txt", d)

	s := box.Bytes("foo.txt")
	r.Equal([]byte("foo!"), s)

	s = box.Bytes("idontexist")
	r.Equal([]byte(""), s)
}

func Test_Box_Bytes_Miss(t *testing.T) {
	r := require.New(t)

	box := NewBox(filepath.Join("_fixtures", "templates"))

	s := box.Bytes("foo.txt")
	r.Equal([]byte("FOO!!!\n"), s)

	s = box.Bytes("idontexist")
	r.Equal([]byte(""), s)
}

func Test_Box_MustBytes(t *testing.T) {
	r := require.New(t)

	box := NewBox("./templates")
	d := resolver.NewInMemory(map[string]file.File{
		"foo.txt": file.NewFile("foo.txt", []byte("foo!")),
	})
	box.SetResolver("foo.txt", d)

	s, err := box.MustBytes("foo.txt")
	r.NoError(err)
	r.Equal("foo!", string(s))

	s, err = box.MustBytes("idontexist")
	r.Error(err)
	r.Equal("", string(s))
}

func Test_Box_MustBytes_Miss(t *testing.T) {
	r := require.New(t)

	box := NewBox("./_fixtures/templates")

	s, err := box.MustBytes("foo.txt")
	r.NoError(err)
	r.Equal("FOO!!!\n", string(s))

	s, err = box.MustBytes("idontexist")
	r.Error(err)
	r.Equal("", string(s))
}

func Test_Box_Has(t *testing.T) {
	r := require.New(t)

	box := NewBox("./templates")
	d := resolver.NewInMemory(map[string]file.File{
		"foo.txt": file.NewFile("foo.txt", []byte("foo!")),
	})
	box.SetResolver("foo.txt", d)

	r.True(box.Has("foo.txt"))
	r.False(box.Has("idontexist"))
}

func Test_Box_Open(t *testing.T) {
	r := require.New(t)

	d := resolver.NewInMemory(map[string]file.File{
		"foo.txt": file.NewFile("foo.txt", []byte("foo!")),
	})
	box := NewBox("./templates")

	box.SetResolver("foo.txt", d)

	f, err := box.Open("foo.txt")
	r.NoError(err)
	r.NotZero(f)

	f, err = box.Open("idontexist")
	r.Error(err)
	r.Zero(f)
}

func Test_Box_List(t *testing.T) {
	r := require.New(t)

	box := NewBox(filepath.Join("_fixtures", "list_test"))
	r.NoError(box.AddString(filepath.Join("d", "d.txt"), "D"))

	act := box.List()
	exp := []string{"a.txt", filepath.Join("b", "b.txt"), filepath.Join("b", "b2.txt"), filepath.Join("c", "c.txt"), filepath.Join("d", "d.txt")}
	r.Equal(exp, act)
}
