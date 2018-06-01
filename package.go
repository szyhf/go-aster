package aster

import (
	"fmt"
	"go/ast"
	"strings"
)

type PackageType struct {
	Name       string           `json:",omitempty"`
	Imports    []*ImportType    `json:",omitempty"`
	Interfaces []*InterfaceType `json:",omitempty"`
	Structs    []*StructType    `json:",omitempty"`
	Funcs      []*FuncType      `json:",omitempty"` // 不考虑Method
	Methods    []*MethodType    `json:",omitempty"`

	importSet map[string]struct{}
}

func NewPackageType(pkg *ast.Package) (*PackageType, error) {
	var err error
	pkgTyp := &PackageType{
		Name:    pkg.Name,
		Imports: make([]*ImportType, 0, 16),
		Structs: make([]*StructType, 0, 64),
		Funcs:   make([]*FuncType, 0, 64),

		importSet: make(map[string]struct{}, 32),
	}
	for _, astFile := range pkg.Files {
		for _, decl := range astFile.Decls {
			switch curDecl := decl.(type) {
			case *ast.GenDecl:
				err = pkgTyp.ParseGenDecl(curDecl)
			case *ast.FuncDecl:
				err = pkgTyp.ParseFuncDecl(curDecl)
			case *ast.BadDecl:

			}
			if err != nil {
				return nil, err
			}
		}
	}

	// 把Method统计到对应的Struct中
	structMap := make(map[string]*StructType, len(pkgTyp.Structs))
	for _, structType := range pkgTyp.Structs {
		structMap[structType.Name] = structType
	}
	for _, methodType := range pkgTyp.Methods {
		receiverName := strings.TrimLeft(methodType.Receiver.Type.String(), "*")
		structType, ok := structMap[receiverName]
		if !ok {
			return nil, fmt.Errorf("NewPackageType: unreslove method receiver.Name = %s", receiverName)
		}
		structType.Methods = append(structType.Methods, methodType)
	}

	return pkgTyp, err
}

func (this *PackageType) ParseFuncDecl(funcDecl *ast.FuncDecl) error {
	if funcDecl.Recv == nil {
		funcType, err := NewFuncTypeByASTDecl(funcDecl)
		if err != nil {
			return err
		}
		this.Funcs = append(this.Funcs, funcType)
	} else {
		methodType, err := NewMethodType(funcDecl)
		if err != nil {
			return err
		}
		this.Methods = append(this.Methods, methodType)
	}
	return nil
}

func (this *PackageType) ParseGenDecl(genDecl *ast.GenDecl) error {
	for _, spec := range genDecl.Specs {
		var err error
		switch curSpec := spec.(type) {
		case *ast.TypeSpec:
			// 以type声明的顶级变量
			err = this.ParseTypeSpec(genDecl, curSpec)
		case *ast.ImportSpec:
			// 以import声明的顶级变量
			err = this.ParseImportSpec(genDecl, curSpec)
		case *ast.ValueSpec:
			// 以var声明的顶级变量
			err = this.ParseValueSpec(genDecl, curSpec)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *PackageType) ParseTypeSpec(astGenDecl *ast.GenDecl, typeSpec *ast.TypeSpec) error {
	switch typeExpr := typeSpec.Type.(type) {
	case *ast.StructType:
		// 是type struct
		structType, err := NewStructType(astGenDecl, typeSpec, typeExpr)
		if err != nil {
			return err
		}
		this.Structs = append(this.Structs, structType)
	case *ast.InterfaceType:
		interfaceType, err := NewInterfaceType(astGenDecl, typeSpec, typeExpr)
		if err != nil {
			return err
		}
		this.Interfaces = append(this.Interfaces, interfaceType)
	case *ast.ArrayType:
	case *ast.FuncType:
	case *ast.MapType:
	case *ast.ChanType:
	case *ast.Ident:
		// fmt.Println(typeExpr.Name)
	default:
		return fmt.Errorf("PackageType.ParseTypeSpec()未处理的*ast.TypeSpec=%T", typeExpr)
	}
	return nil
}

func (this *PackageType) ParseImportSpec(astGenDecl *ast.GenDecl, importSpec *ast.ImportSpec) error {
	if _, ok := this.importSet[importSpec.Path.Value]; !ok {
		this.Imports = append(this.Imports, &ImportType{
			Alias: func() string {
				if importSpec.Name != nil {
					return importSpec.Name.Name
				}
				return ""
			}(),
			// delete the '"' prefix and sufix
			Name: strings.Trim(importSpec.Path.Value, `"`),
		})
		this.importSet[importSpec.Path.Value] = struct{}{}
	}
	return nil
}

func (this *PackageType) ParseValueSpec(astGenDecl *ast.GenDecl, valueSpec *ast.ValueSpec) error {
	return nil
}

func (this *PackageType) String() string {
	sb := strings.Builder{}
	sb.WriteString("package " + this.Name + "\n\n")
	sb.WriteString("import (\n")
	for _, importType := range this.Imports {
		if importType.Alias != "" {
			sb.WriteString("\t" + importType.Alias + ` "` + importType.Name + `"` + "\n")
		} else {
			sb.WriteString("\t" + ` "` + importType.Name + `"` + "\n")
		}
	}
	sb.WriteString(")\n\n")

	for _, interfaceType := range this.Interfaces {
		sb.WriteString(interfaceType.String() + "\n")
	}

	for _, structType := range this.Structs {
		sb.WriteString(structType.String() + "\n")
	}

	for _, methodType := range this.Methods {
		sb.WriteString(methodType.String() + "\n")
	}

	for _, funcType := range this.Funcs {
		sb.WriteString(funcType.String() + "\n")
	}

	return sb.String()
}
