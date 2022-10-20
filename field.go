package aster

import (
	"fmt"
	"go/ast"
	"strings"
)

// `go/ast`包中的`Field`结构的映射
// 用于描述一个字段，可能是结构、接口中的字段，也可能是函数的参数、方法的接受者、全局变量等
type FieldType struct {
	Name string    `json:",omitempty"`
	Docs []Comment `json:",omitempty"`
	Type *TypeType `json:",omitempty"`
}

// 因为语法上存在通过省略类型而实际有多个参数的情况，所以返回值是数组（例如`x,y string`）
// 没有error的情况下，返回的数组保底有一个
func NewFieldTypes(astField *ast.Field) ([]*FieldType, error) {
	astNames := make([]*ast.Ident, len(astField.Names))
	copy(astNames, astField.Names)
	if len(astNames) == 0 {
		// 填充一个匿名对象便于执行业务
		astNames = append(astNames, &ast.Ident{})
	}
	resFieldTypes := make([]*FieldType, len(astNames))
	// if len(astNames) > 1 {
	// 	log.Printf("%+v", astField.Names)
	// 	defer func() {
	// 		log.Printf("%+v", resFieldTypes)
	// 	}()
	// }
	for i, astName := range astNames {
		fieldType := &FieldType{}
		typeType, err := NewTypeType(astField.Type)
		if err != nil {
			return nil, err
		}
		fieldType.Type = typeType

		if astField.Doc != nil {
			fieldType.Docs = make([]Comment, 0, len(astField.Doc.List))
			for _, docComment := range astField.Doc.List {
				fieldType.Docs = append(fieldType.Docs, docComment.Text)
			}
		}
		fieldType.Name = astName.Name

		resFieldTypes[i] = fieldType
	}

	return resFieldTypes, nil
}

// 如果是匿名结构体，则直接返回原始名，形如`FiledType`
// 如果是有字段名的类型，则返回`Field FiledType`
func (this *FieldType) GetDecl() string {
	if this.Name == "" {
		return this.Type.GetDecl()
	}
	return fmt.Sprintf("%s %s", this.Name, this.Type.GetDecl())
}

type StructFieldType struct {
	FieldType
	Tag TagType `json:",omitempty"`
}

// 因为语法上存在通过省略类型而实际有多个参数的情况，所以返回值是数组（例如`x,y string`）
// 没有error的情况下，返回的数组保底有一个
func NewStructFieldType(astField *ast.Field) ([]*StructFieldType, error) {
	fieldTypes, err := NewFieldTypes(astField)
	if err != nil {
		return nil, err
	}
	resStructFieldTypes := make([]*StructFieldType, len(fieldTypes))
	for i, fieldType := range fieldTypes {
		structFieldType := &StructFieldType{
			FieldType: *fieldType,
		}
		if astField.Tag != nil {
			structFieldType.Tag = TagType(strings.Trim(astField.Tag.Value, "`"))
		}
		resStructFieldTypes[i] = structFieldType
	}

	return resStructFieldTypes, nil
}

// 如果是匿名结构体，则直接返回原始名，形如`FiledType Tag`
// 如果是有字段名的类型，则返回`Field FiledType Tag`
func (this *StructFieldType) GetDecl() string {
	return fmt.Sprintf("%s %s", this.FieldType.GetDecl(), this.Tag)
}
