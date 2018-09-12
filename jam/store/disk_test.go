package store

import (
	"path/filepath"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/packr/jam/parser"
	"github.com/stretchr/testify/require"
)

func Test_Disk_FileNames(t *testing.T) {
	r := require.New(t)

	d := &Disk{}

	box := parser.NewBox("Test_Disk_FileNames", "./_fixtures/disk/franklin")
	names, err := d.FileNames(box)
	r.NoError(err)
	r.Len(names, 2)

	r.Equal("aretha.txt", filepath.Base(names[0]))
	r.Equal("think.txt", filepath.Base(names[1]))
}

func Test_Disk_Files(t *testing.T) {
	r := require.New(t)

	d := &Disk{}

	box := parser.NewBox("Test_Disk_Files", "./_fixtures/disk/franklin")
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

func Test_Disk_Pack(t *testing.T) {
	r := require.New(t)

	d := NewDisk("", "")

	p, err := parser.NewFromRoots([]string{"./_fixtures/disk-pack"})
	r.NoError(err)
	boxes, err := p.Run()
	r.NoError(err)

	for _, b := range boxes {
		r.NoError(d.Pack(b))
	}

	global := d.global
	r.Len(global, 3)

	r.Len(d.boxes, 3)

}

func Test_Disk_Packed_Test(t *testing.T) {
	r := require.New(t)

	b := packr.NewBox("simpsons")

	s, err := b.MustString("parents/homer.txt")
	r.NoError(err)
	r.Equal("HOMER Simpson\n", s)

	s, err = b.MustString("parents/marge.txt")
	r.NoError(err)
	r.Equal("MARGE Simpson\n", s)

	_, err = b.MustString("idontexist")
	r.Error(err)
}

func Test_Disk_Close(t *testing.T) {
	r := require.New(t)

	p, err := parser.NewFromRoots([]string{"./_fixtures/disk-pack"})
	r.NoError(err)
	boxes, err := p.Run()
	r.NoError(err)

	d := NewDisk("", "")
	for _, b := range boxes {
		r.NoError(d.Pack(b))
	}
	r.NoError(d.Close())
}
