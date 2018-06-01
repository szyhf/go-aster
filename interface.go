package aster

import (
	"fmt"
	"go/ast"
	"strings"
)

type InterfaceType struct {
	Name  string `json:",omitempty"`
	Funcs []*InterfaceFuncType
	Docs  []Comment
}

func NewInterfaceType(astGenDecl *ast.GenDecl, typeSpec *ast.TypeSpec, astInterface *ast.InterfaceType) (*InterfaceType, error) {
	interfaceType := &InterfaceType{
		Name: typeSpec.Name.String(),
	}

	if astInterface.Methods != nil {
		interfaceType.Funcs = make([]*InterfaceFuncType, 0, astInterface.Methods.NumFields())
		for _, methodField := range astInterface.Methods.List {
			switch astExpr := methodField.Type.(type) {
			case *ast.FuncType:
				err := interfaceType.ParseInterfaceFuncType(methodField, astExpr)
				if err != nil {
					return nil, err
				}
			case *ast.Ident:
				// 在interface中嵌入的子interface
				// TODO: 会导致无法拿到Interface的所有方法，需要额外处理。
				// fmt.Println("NewInterfaceType(): 在interface中嵌入的子interface，会导致无法拿到Interface的所有方法，需要额外处理。")
			default:
				return nil, fmt.Errorf("NewInterfaceType()未处理的MethodField: %T", astExpr)
			}

		}
	}

	if astGenDecl.Doc != nil {
		interfaceType.Docs = make([]Comment, 0, len(astGenDecl.Doc.List))
		for _, doc := range astGenDecl.Doc.List {
			interfaceType.Docs = append(interfaceType.Docs, doc.Text)
		}
	}

	return interfaceType, nil
}

func (this *InterfaceType) ParseInterfaceFuncType(astField *ast.Field, astFuncType *ast.FuncType) error {
	funcType, err := NewInterfaceFuncType(astField, astFuncType)
	if err == nil {
		this.Funcs = append(this.Funcs, funcType)
	}
	return err
}

func (this *InterfaceType) String() string {
	sb := strings.Builder{}

	for _, doc := range this.Docs {
		sb.WriteString(doc)
	}

	sb.WriteString("type " + this.Name + " interface {\n")
	for _, fun := range this.Funcs {
		sb.WriteString("\t" + fun.Name + "(")
		for i, paramType := range fun.Params {
			sb.WriteString(paramType.String())
			if i < len(fun.Params)-1 {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(") ")

		if len(fun.Results) > 1 {
			sb.WriteString("(")
		}
		for i, resultType := range fun.Results {
			sb.WriteString(resultType.String())
			if i < len(fun.Results)-1 {
				sb.WriteString(", ")
			}
		}
		if len(fun.Results) > 1 {
			sb.WriteString(")")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}
