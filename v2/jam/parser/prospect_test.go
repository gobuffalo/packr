package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_IsProspect(t *testing.T) {
	table := []struct {
		root string
		path string
		pass bool
	}{
		{"", "foo/.git/config", false},
		{"", "foo/.git/baz.go", false},
		{"", "a.go", true},
		{"", ".", true},
		{"", "a/b.go", true},
		{"", "a/b_test.go", false},
		{"", "a/b-packr.go", false},
		{"", "a/vendor/b.go", false},
		{"", "a/_c/c.go", false},
		{"", "a/_c/e/fe/f/c.go", false},
		{"", "a/d/_d.go", false},
		{"", "a/d/", false},
		{".", ".", true},
		{"a", "a/b.go", true},
		{"a/vendor", "a/vendor/b.go", true},
		{"a", "a/vendor/b.go", false},
		{".ci", ".ci/a/b.go", true},
		{"a", "a/.ci/b.go", false},
	}

	for _, tt := range table {
		t.Run(fmt.Sprintf("%s:%s", tt.root, tt.path), func(st *testing.T) {
			r := require.New(st)
			if tt.pass {
				r.True(IsProspect(tt.root, tt.path, ".", "_"))
			} else {
				r.False(IsProspect(tt.root, tt.path, ".", "_"))
			}
		})
	}
}
