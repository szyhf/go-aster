package aster

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func ParseDir(dirPath string, fileFilter func(os.FileInfo) bool) ([]*PackageType, error) {
	var err error
	fSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fSet, dirPath, fileFilter, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	pkgsTyp := make([]*PackageType, 0, len(pkgs))
	for _, astPkg := range pkgs {
		pkgTyp, err := NewPackageType(astPkg)
		if err != nil {
			return nil, err
		}
		pkgsTyp = append(pkgsTyp, pkgTyp)
	}
	return pkgsTyp, nil
}

func ParseFile(filePath string) (*PackageType, error) {
	fSet := token.NewFileSet()
	var astPkg *ast.Package
	if src, err := parser.ParseFile(fSet, filePath, nil, 0); err == nil {
		astPkg = &ast.Package{
			Name:  src.Name.Name,
			Files: make(map[string]*ast.File),
		}
		astPkg.Files[filePath] = src
		pkgTyp, err := NewPackageType(astPkg)
		if err != nil {
			return nil, err
		}
		return pkgTyp, nil
	} else {
		return nil, err
	}
}
