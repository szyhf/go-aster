package aster

import (
	"go/ast"
	"strings"
)

type StructType struct {
	// 当前结构所属的PackageType的引用
	PackageType *PackageType

	Name       string             `json:",omitempty"`
	TypeParams []*StructFieldType `json:",omitempty"`
	Fields     []*StructFieldType `json:",omitempty"`
	Methods    []*MethodType      `json:",omitempty"`
	Docs       []Comment          `json:",omitempty"`
}

func (pkgType *PackageType) NewStructType(astGenDecl *ast.GenDecl, typeSpec *ast.TypeSpec, astStructType *ast.StructType) (*StructType, error) {
	structType := &StructType{
		PackageType: pkgType,

		Name:    typeSpec.Name.Name,
		Methods: make([]*MethodType, 0, 16),
	}

	if typeSpec.TypeParams != nil {
		structType.TypeParams = make([]*StructFieldType, 0, len(typeSpec.TypeParams.List))
		for _, astField := range typeSpec.TypeParams.List {
			err := structType.ParseTypeParam(astField)
			if err != nil {
				return nil, err
			}
		}
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

func (this *StructType) ParseTypeParam(astField *ast.Field) error {
	fieldTypes, err := NewStructFieldType(astField)
	if err != nil {
		return err
	}
	// log.Println(gjson.MustEncodeString(fieldTypes))
	this.TypeParams = append(this.TypeParams, fieldTypes...)
	return nil
}

func (this *StructType) ParseStructField(astField *ast.Field) error {
	fieldTypes, err := NewStructFieldType(astField)
	if err != nil {
		return err
	}
	this.Fields = append(this.Fields, fieldTypes...)
	return nil
}

// 声明时的名字，形如 Struct[T any]
func (this *StructType) GetDeclName() string {
	return this.Name + this.getTypeParamsDeclName()
}

// 作为接受者时的名字，形如 Struct[T]
func (this *StructType) GetRecvName() string {
	return this.Name + this.getTypeParamsRecvName()
}

// 获取完整的结构体声明
func (this *StructType) GetDecl() string {
	sb := strings.Builder{}
	sb.WriteString("type " + this.GetDeclName() + " struct {\n")
	for _, field := range this.Fields {
		sb.WriteString("\t" + field.GetDecl() + "\n")
	}
	sb.WriteString("}\n")
	return sb.String()
}

func (this *StructType) getTypeParamsDeclName() string {
	if len(this.TypeParams) == 0 {
		return ""
	}
	sb := &strings.Builder{}
	sb.WriteString("[")
	sb.WriteString(this.TypeParams[0].FieldType.GetDecl())
	for _, indiceTyp := range this.TypeParams[1:] {
		sb.WriteString(",")
		sb.WriteString(strings.TrimSpace(indiceTyp.FieldType.GetDecl()))
	}
	sb.WriteString("]")
	return sb.String()
}

func (this *StructType) getTypeParamsRecvName() string {
	if len(this.TypeParams) == 0 {
		return ""
	}
	sb := &strings.Builder{}
	sb.WriteString("[")
	sb.WriteString(this.TypeParams[0].FieldType.Name)
	for _, indiceTyp := range this.TypeParams[1:] {
		sb.WriteString(",")
		sb.WriteString(strings.TrimSpace(indiceTyp.FieldType.Name))
	}
	sb.WriteString("]")
	return sb.String()
}
