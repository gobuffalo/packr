package parser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/gobuffalo/packr/v2/plog"
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
		plog.Debug(p, "Run", "parsing", pros.Name())
		v := NewVisitor(pros)
		pbr, err := v.Run()
		if err != nil {
			return boxes, errors.WithStack(err)
		}
		for _, b := range pbr {
			plog.Debug(p, "Run", "file", pros.Name(), "box", b.Name)
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

func (r RootsOptions) String() string {
	x, _ := json.Marshal(r)
	return string(x)
}

// NewFromRoots scans the file roots provided and returns a
// new Parser containing the prospects
func NewFromRoots(roots []string, opts *RootsOptions) (*Parser, error) {
	if opts == nil {
		opts = &RootsOptions{}
	}
	p := New()
	plog.Debug(p, "NewFromRoots", "roots", roots, "options", opts)
	callback := func(path string, de *godirwalk.Dirent) error {
		if IsProspect(path, opts.Ignores...) && de.IsDir() {
			roots = append(roots, path)
			return nil
		}
		if de.IsDir() {
			return filepath.SkipDir
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
	fd := &finder{}
	for _, r := range roots {
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
