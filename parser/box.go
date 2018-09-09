package parser

import "encoding/json"

type Box struct {
	Name       string
	Path       string
	Package    string
	PWD        string
	PackageDir string
}

func (b Box) String() string {
	x, _ := json.Marshal(b)
	return string(x)
}
