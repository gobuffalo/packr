package store

import (
	"io/ioutil"
	"testing"

	"github.com/gobuffalo/packr/file"
	"github.com/stretchr/testify/require"
)

func Test_inMemory(t *testing.T) {
	r := require.New(t)

	s := NewInMemory()

	r.NoError(s.Pack("foo.txt", file.NewFile("foo.txt", []byte("foo!"))))

	fm, ok := s.(file.FileMappable)
	r.True(ok)
	m := fm.FileMap()
	r.Len(m, 1)

	f, ok := m["foo.txt"]
	r.True(ok)
	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.Equal("foo!", string(b))
}
