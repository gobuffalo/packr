package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"strings"

	"github.com/gobuffalo/packd"
	"github.com/pkg/errors"
)

// ParsedFile ...
type ParsedFile struct {
	File    packd.SimpleFile
	FileSet *token.FileSet
	Ast     *ast.File
	Lines   []string
}

// ParseFileMode ...
func ParseFileMode(gf packd.SimpleFile, mode parser.Mode) (ParsedFile, error) {
	pf := ParsedFile{
		FileSet: token.NewFileSet(),
		File:    gf,
	}

	src := gf.String()
	f, err := parser.ParseFile(pf.FileSet, gf.Name(), src, mode)
	if err != nil && errors.Cause(err) != io.EOF {
		return pf, errors.WithStack(err)
	}
	pf.Ast = f

	pf.Lines = strings.Split(src, "\n")
	return pf, nil
}

// ParseFile ...
func ParseFile(gf packd.SimpleFile) (ParsedFile, error) {
	return ParseFileMode(gf, 0)
}
