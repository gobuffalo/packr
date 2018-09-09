package store

import (
	"path/filepath"
	"testing"

	"github.com/gobuffalo/packr/costello/parser"
	"github.com/stretchr/testify/require"
)

func Test_Disk_FileNames(t *testing.T) {
	r := require.New(t)

	d := &Disk{}

	box := parser.NewBox("Test_Disk_FileNames", "../store/_fixtures/disk/franklin")
	names, err := d.FileNames(box)
	r.NoError(err)
	r.Len(names, 2)

	r.Equal("aretha.txt", filepath.Base(names[0]))
	r.Equal("think.txt", filepath.Base(names[1]))
}

func Test_Disk_Files(t *testing.T) {
	r := require.New(t)

	d := &Disk{}

	box := parser.NewBox("Test_Disk_Files", "../store/_fixtures/disk/franklin")
	files, err := d.Files(box)
	r.NoError(err)
	r.Len(files, 2)

	f := files[0]
	r.Equal("aretha.txt", filepath.Base(f.Name()))
	r.Equal("RESPECT!\n", f.String())

	f = files[1]
	r.Equal("think.txt", filepath.Base(f.Name()))
	r.Equal("THINK!\n", f.String())
}
