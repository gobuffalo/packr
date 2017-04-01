package packr

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Box_String(t *testing.T) {
	r := require.New(t)
	s := testBox.String("hello.txt")
	// Need to compare the file, that test works on different machines
	testFile, err := ioutil.ReadFile(filepath.Join("fixtures", "hello.txt"))
	if err != nil {
		t.Error("Cannot read test file: ", err)
	}
	r.Equal(string(testFile), s)
}

func Test_Box_MustString(t *testing.T) {
	r := require.New(t)
	_, err := testBox.MustString("idontexist.txt")
	r.Error(err)
}

func Test_Box_Bytes(t *testing.T) {
	r := require.New(t)
	s := testBox.Bytes("hello.txt")
	// Need to compare the file, that test works on different machines
	testFile, err := ioutil.ReadFile(filepath.Join("fixtures", "hello.txt"))
	if err != nil {
		t.Error("Cannot read test file: ", err)
	}
	r.Equal(testFile, s)
}

func Test_Box_MustBytes(t *testing.T) {
	r := require.New(t)
	_, err := testBox.MustBytes("idontexist.txt")
	r.Error(err)
}

func Test_Box_Walk_Physical(t *testing.T) {
	r := require.New(t)
	count := 0
	err := testBox.Walk(func(path string, f File) error {
		count++
		return nil
	})
	r.NoError(err)
	r.Equal(2, count)
}

func Test_Box_Walk_Virtual(t *testing.T) {
	r := require.New(t)
	count := 0
	err := virtualBox.Walk(func(path string, f File) error {
		count++
		return nil
	})
	r.NoError(err)
	r.Equal(3, count)
}
