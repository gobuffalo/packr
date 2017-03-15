package packr

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Box_String(t *testing.T) {
	r := require.New(t)
	s := testBox.String("hello.txt")
	r.Equal("hello world!\n", s)
}

func Test_Box_MustString(t *testing.T) {
	r := require.New(t)
	_, err := testBox.MustString("idontexist.txt")
	r.Error(err)
}

func Test_Box_Bytes(t *testing.T) {
	r := require.New(t)
	s := testBox.Bytes("hello.txt")
	r.Equal([]byte("hello world!\n"), s)
}

func Test_Box_MustBytes(t *testing.T) {
	r := require.New(t)
	_, err := testBox.MustBytes("idontexist.txt")
	r.Error(err)
}
