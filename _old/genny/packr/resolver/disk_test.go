package resolver

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DiskResolver(t *testing.T) {
	r := require.New(t)

	dr := &DiskResolver{
		Roots: []string{filepath.Join("fixtures", "good")},
	}

	r.NoError(dr.Resolve())

	pros := dr.Prospects()
	r.Len(pros, 1)

	f := pros[0]
	r.Equal(filepath.Join("fixtures", "good", "main.go"), f.Name())

	boxes := dr.Boxes()
	r.Len(boxes, 1)

	box, ok := boxes[dr.Roots[0]]
	r.True(ok)
	r.Len(box, 2)
}
