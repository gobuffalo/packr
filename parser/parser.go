package parser

import (
	"fmt"

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
