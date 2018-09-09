package parser

import (
	"fmt"

	"github.com/pkg/errors"
)

type Parser struct {
	Prospects []*File
}

func (p *Parser) Run() error {
	for _, pros := range p.Prospects {
		fmt.Println("Parser: parsing", pros.Name())
		v := NewVisitor(pros)
		boxes, err := v.Run()
		if err != nil {
			return errors.WithStack(err)
		}
		for _, n := range boxes {
			fmt.Println("### n ->", n)
		}
	}
	return nil
}
