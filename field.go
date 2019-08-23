package aster

import (
	"fmt"
	"go/ast"
	"strings"
)

type FieldType struct {
	Name string    `json:",omitempty"`
	Docs []Comment `json:",omitempty"`
	Type *TypeType `json:",omitempty"`
}

func NewFieldType(astField *ast.Field) (*FieldType, error) {
	fieldType := &FieldType{}

	typeType, err := NewTypeType(astField.Type)
	if err != nil {
		return nil, err
	}
	fieldType.Type = typeType

	if len(astField.Names) > 0 {
		// 字段名可能是匿名
		fieldType.Name = astField.Names[0].Name
	}
	if astField.Doc != nil {
		fieldType.Docs = make([]Comment, 0, len(astField.Doc.List))
		for _, docComment := range astField.Doc.List {
			fieldType.Docs = append(fieldType.Docs, docComment.Text)
		}
	}

	return fieldType, nil
}

func (this *FieldType) String() string {
	if this.Name == "" {
		return this.Type.String()
	}
	return fmt.Sprintf("%s %s", this.Name, this.Type)
}

type StructFieldType struct {
	FieldType
	Tag TagType `json:",omitempty"`
}

func NewStructFieldType(astField *ast.Field) (*StructFieldType, error) {
	fieldType, err := NewFieldType(astField)
	if err != nil {
		return nil, err
	}
	structFieldType := &StructFieldType{
		FieldType: *fieldType,
	}
	if astField.Tag != nil {
		structFieldType.Tag = TagType(strings.Trim(astField.Tag.Value, "`"))
	}
	return structFieldType, nil
}

func (this *StructFieldType) String() string {
	return fmt.Sprintf("%s %s", this.FieldType.String(), this.Tag)
}
