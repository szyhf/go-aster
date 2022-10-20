package aster

import (
	"go/ast"
	"strings"
)

type MethodType struct {
	FuncType
	Receiver *FieldType `json:",omitemtpy"`

	FuncDecl *ast.FuncDecl `json:"-"`
}

func NewMethodType(astFuncDecl *ast.FuncDecl) (*MethodType, error) {
	funcType, err := NewFuncTypeByASTDecl(astFuncDecl)
	if err != nil {
		return nil, err
	}

	recvier := astFuncDecl.Recv.List[0]
	fieldTypes, err := NewFieldTypes(recvier)
	if err != nil {
		return nil, err
	}
	// recevier没有省略语法的场景
	fieldType := fieldTypes[0]

	methodType := &MethodType{
		FuncType: *funcType,
		FuncDecl: astFuncDecl,
		Receiver: fieldType,
	}

	return methodType, nil
}

func (this *MethodType) String() string {
	sb := strings.Builder{}

	sb.WriteString("func (" + this.Receiver.Name + " " + this.Receiver.Type.GetDecl() + ") " + this.Name + "(")

	for i, paramType := range this.Params {
		sb.WriteString(paramType.GetDecl())
		if i < len(this.Params)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(") ")

	if len(this.Results) > 1 {
		sb.WriteString("(")
	}
	for i, resultType := range this.Results {
		sb.WriteString(resultType.GetDecl())
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
