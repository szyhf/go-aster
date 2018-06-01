package aster

import (
	"fmt"
	"go/ast"
	"strings"
)

type FuncType struct {
	Name    string       `json:",omitempty"`
	Params  []*FieldType `json:",omitempty"`
	Results []*FieldType `json:",omitempty"`

	astBolckStmt *ast.BlockStmt
}

func NewFuncTypeByASTField(astField *ast.Field, astFuncType *ast.FuncType) (*FuncType, error) {
	funcType := &FuncType{}

	if len(astField.Names) > 0 {
		funcType.Name = astField.Names[0].Name
	}
	if astFuncType.Params != nil {
		funcType.Params = make([]*FieldType, 0, astFuncType.Params.NumFields())
		for _, astParamField := range astFuncType.Params.List {
			err := funcType.ParseParam(astParamField)
			if err != nil {
				return nil, err
			}
		}
	}
	if astFuncType.Results != nil {
		funcType.Results = make([]*FieldType, 0, astFuncType.Results.NumFields())
		for _, astResultField := range astFuncType.Results.List {
			err := funcType.ParseResult(astResultField)
			if err != nil {
				return nil, err
			}
		}
	}
	return funcType, nil
}

func NewFuncTypeByASTDecl(astDecl *ast.FuncDecl) (*FuncType, error) {
	funcType := &FuncType{}
	astFuncType := astDecl.Type

	if astDecl.Name != nil {
		funcType.Name = astDecl.Name.Name
	}
	if astFuncType.Params != nil {
		funcType.Params = make([]*FieldType, 0, astFuncType.Params.NumFields())
		for _, astParamField := range astFuncType.Params.List {
			err := funcType.ParseParam(astParamField)
			if err != nil {
				return nil, err
			}
		}
	}
	if astFuncType.Results != nil {
		funcType.Results = make([]*FieldType, 0, astFuncType.Results.NumFields())
		for _, astResultField := range astFuncType.Results.List {
			err := funcType.ParseResult(astResultField)
			if err != nil {
				return nil, err
			}
		}
	}

	funcType.astBolckStmt = astDecl.Body
	return funcType, nil
}

func (this *FuncType) ParseParam(astParamField *ast.Field) error {
	field, err := NewFieldType(astParamField)
	if err == nil {
		this.Params = append(this.Params, field)
	}
	return err
}

func (this *FuncType) ParseResult(astResultField *ast.Field) error {
	field, err := NewFieldType(astResultField)
	if err == nil {
		this.Results = append(this.Results, field)
	}
	return err
}

// 不是所有FuncType都有body，可能只是个声明。
func (this *FuncType) GetASTBlockSTMT() *ast.BlockStmt {
	return this.astBolckStmt
}

func (this *FuncType) String() string {

	sb := strings.Builder{}

	sb.WriteString("func " + this.Name + "(")

	for i, paramType := range this.Params {
		sb.WriteString(paramType.String())
		if i < len(this.Params)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(") ")

	if len(this.Results) > 1 {
		sb.WriteString("(")
	}
	for i, resultType := range this.Results {
		sb.WriteString(resultType.String())
		if i < len(this.Results)-1 {
			sb.WriteString(", ")
		}
	}
	if len(this.Results) > 1 {
		sb.WriteString(")")
	}
	sb.WriteString(" {\n\t// ...\n}\n")
	return sb.String()
}

type InterfaceFuncType struct {
	FuncType

	ASTField    *ast.Field
	ASTFuncType *ast.FuncType
}

func NewInterfaceFuncType(astField *ast.Field, astFuncType *ast.FuncType) (*InterfaceFuncType, error) {
	funcType := &InterfaceFuncType{
		ASTField:    astField,
		ASTFuncType: astFuncType,
	}

	if len(astField.Names) > 0 {
		funcType.Name = astField.Names[0].Name
	}
	if astFuncType.Params != nil {
		funcType.Params = make([]*FieldType, 0, astFuncType.Params.NumFields())
		for _, astParamField := range astFuncType.Params.List {
			err := funcType.ParseParam(astParamField)
			if err != nil {
				return nil, err
			}
		}
	}
	if astFuncType.Results != nil {
		funcType.Results = make([]*FieldType, 0, astFuncType.Results.NumFields())
		for _, astResultField := range astFuncType.Results.List {
			err := funcType.ParseResult(astResultField)
			if err != nil {
				return nil, err
			}
		}
	}
	return funcType, nil
}

func (this *InterfaceFuncType) String() string {
	return fmt.Sprintf("InterfaceFuncType.String()")
}
