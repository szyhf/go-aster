package aster

import (
	"go/ast"
	"strings"
)

type StructType struct {
	Name    string             `json:",omitempty"`
	Fields  []*StructFieldType `json:",omitempty"`
	Methods []*MethodType      `json:",omitempty"`
	Docs    []Comment          `json:",omitempty"`
}

func NewStructType(astGenDecl *ast.GenDecl, typeSpec *ast.TypeSpec, astStructType *ast.StructType) (*StructType, error) {
	structType := &StructType{
		Name:    typeSpec.Name.Name,
		Methods: make([]*MethodType, 0, 16),
	}

	if astStructType.Fields != nil {
		structType.Fields = make([]*StructFieldType, 0, astStructType.Fields.NumFields())
		for _, astField := range astStructType.Fields.List {
			err := structType.ParseStructField(astField)
			if err != nil {
				return nil, err
			}
		}
	}

	if astGenDecl.Doc != nil {
		structType.Docs = make([]Comment, 0, len(astGenDecl.Doc.List))
		for _, doc := range astGenDecl.Doc.List {
			structType.Docs = append(structType.Docs, doc.Text)
		}
	}
	return structType, nil
}

func (this *StructType) ParseStructField(astField *ast.Field) error {
	fieldType, err := NewStructFieldType(astField)
	if err != nil {
		return err
	}
	this.Fields = append(this.Fields, fieldType)
	return nil
}

func (this *StructType) String() string {
	sb := strings.Builder{}
	sb.WriteString("type " + this.Name + " struct {\n")
	for _, field := range this.Fields {
		sb.WriteString("\t" + field.String() + "\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}
