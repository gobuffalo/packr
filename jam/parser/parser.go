package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

// Parser to find boxes
type Parser struct {
	Prospects     []*File // a list of files to check for boxes
	IgnoreImports bool
}

// Run the parser and run any boxes found
func (p *Parser) Run() ([]*Box, error) {
	var boxes []*Box
	for _, pros := range p.Prospects {
		// fmt.Println("Parser: parsing", pros.Name())
		v := NewVisitor(pros)
		pbr, err := v.Run()
		if err != nil {
			return boxes, errors.WithStack(err)
		}
		for _, b := range pbr {
			boxes = append(boxes, b)
		}
	}
	return boxes, nil
}

// New Parser from a list of File
func New(prospects ...*File) *Parser {
	return &Parser{
		Prospects: prospects,
	}
}

type RootsOptions struct {
	IgnoreImports bool
	Ignores       []string
}

// NewFromRoots scans the file roots provided and returns a
// new Parser containing the prospects
func NewFromRoots(roots []string, opts *RootsOptions) (*Parser, error) {
	if opts == nil {
		opts = &RootsOptions{}
	}
	fmt.Println("Parser: prospecting roots\n", strings.Join(roots, "\n"))
	p := New()
	callback := func(path string, de *godirwalk.Dirent) error {
		if IsProspect(path, opts.Ignores...) && de.IsDir() {
			roots = append(roots, path)
			return nil
		}
		return nil
	}
	wopts := &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	}
	for _, root := range roots {
		err := godirwalk.Walk(root, wopts)
		if err != nil {
			return p, errors.WithStack(err)
		}
	}

	dd := map[string]string{}
	for _, r := range roots {
		fd := &finder{
			seen: map[string]string{},
		}
		var names []string
		if opts.IgnoreImports {
			names, _ = fd.findAllGoFiles(r)
		} else {
			names, _ = fd.findAllGoFilesImports(r)
		}
		for _, n := range names {
			if IsProspect(n) {
				dd[n] = n
			}
		}
	}
	for path := range dd {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		p.Prospects = append(p.Prospects, NewFile(path, bytes.NewReader(b)))
	}
	return p, nil
}
