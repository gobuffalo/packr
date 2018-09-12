package parser

import (
	"bytes"
	"fmt"
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
		fmt.Println("Parser: parsing", pros.Name())
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
		if !IsProspect(path, ignore...) {
			if de.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if de.IsDir() {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return errors.WithStack(err)
		}
		p.Prospects = append(p.Prospects, NewFile(path, bytes.NewReader(b)))
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
	return p, nil
}
