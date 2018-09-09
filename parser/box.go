package parser

type Box struct {
	Name       string
	Package    string
	PWD        string
	PackageDir string
	Files      map[string]*File
}
