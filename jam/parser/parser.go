package parser

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
)

type Parser struct {
	Prospects []*File
}

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

func New(prospects ...*File) *Parser {
	return &Parser{
		Prospects: prospects,
	}
}

func NewFromRoots(roots []string, ignore ...string) (*Parser, error) {
	fmt.Println("Parser: prospecting roots\n", strings.Join(roots, "\n"))
	p := New()
	callback := func(path string, de *godirwalk.Dirent) error {
		if IsProspect(path, ignore...) && de.IsDir() {
			roots = append(roots, path)
			return nil
		}
		return nil
	}
	opts := &godirwalk.Options{
		FollowSymbolicLinks: true,
		Callback:            callback,
	}
	for _, root := range roots {
		err := godirwalk.Walk(root, opts)
		if err != nil {
			return p, errors.WithStack(err)
		}
	}

	dd := map[string]string{}
	for _, r := range roots {
		fd := &finder{
			seen: map[string]string{},
		}
		names, _ := fd.findAllGoFiles(r)
		for _, n := range names {
			dd[n] = n
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

type finder struct {
	seen map[string]string
}

func (fd *finder) findAllGoFiles(dir string) ([]string, error) {
	ctx := build.Default
	var names []string
	pkg, err := ctx.ImportDir(dir, 0)

	if err != nil {
		if !strings.Contains(err.Error(), "cannot find package") {
			if _, ok := errors.Cause(err).(*build.NoGoError); !ok {
				return names, errors.WithStack(err)
			}
		}
	}

	if pkg.Goroot {
		return names, nil
	}

	if len(pkg.GoFiles) <= 0 {
		return names, nil
	}
	for _, n := range pkg.GoFiles {
		names = append(names, filepath.Join(pkg.Dir, n))
	}
	for _, imp := range pkg.Imports {
		if _, ok := fd.seen[imp]; ok {
			continue
		}
		fd.seen[imp] = imp
		// fmt.Println("resolving package", pkg.ImportPath, pkg.GoFiles)
		for _, d := range ctx.SrcDirs() {
			ip := filepath.Join(d, imp)
			n, err := fd.findAllGoFiles(ip)
			if err != nil {
				return n, nil
			}
			names = append(names, n...)
		}
	}
	return names, nil
}
