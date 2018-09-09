package parser

import (
	"encoding/json"
	"os"
	"strings"
)

type Box struct {
	Name       string
	Path       string
	AbsPath    string
	Package    string
	PWD        string
	PackageDir string
}

func (b Box) String() string {
	x, _ := json.Marshal(b)
	return string(x)
}

func NewBox(name string, path string) *Box {
	if len(name) == 0 {
		name = path
	}
	name = strings.Replace(name, "\"", "", -1)
	pwd, _ := os.Getwd()
	box := &Box{
		Name: name,
		Path: path,
		PWD:  pwd,
	}
	return box
}
