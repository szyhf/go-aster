package aster

import (
	"fmt"
	"go/ast"
	"strings"
)

type TypeType struct {
	Name string `json:",omitempty"`
	Kind Kind   `json:",omitempty"`

	Elem       *TypeType   `json:",omitempty"`
	TypeParams []*TypeType `json:",omitempty"`
}

func NewTypeType(astExpr ast.Expr) (typeType *TypeType, err error) {
	typeType = &TypeType{}
	switch exprType := astExpr.(type) {
	case *ast.SelectorExpr:
		// 引用其他包
		typeType.Kind = Selector
		typeType.Name = exprType.X.(*ast.Ident).Name + "." + exprType.Sel.Name
		return
	case *ast.Ident:
		typeType.Kind = Ident
		typeType.Name = exprType.Name
		return
	case *ast.StructType:
		typeType.Kind = Struct
		typeType.Name = "struct{}"
		return
	case *ast.FuncType:
		typeType.Kind = Func
		typeType.Name = "" // 匿名方法 Block
		return
	case *ast.StarExpr:
		typeType.Kind = Star
		typeType.Elem, err = NewTypeType(exprType.X)
		return
	case *ast.Ellipsis:
		// 使用Slice表示...的情况 省略符
		typeType.Kind = Ellipsis
		typeType.Elem, err = NewTypeType(exprType.Elt)
		return
	case *ast.ArrayType:
		// 数组
		typeType.Kind = Array
		typeType.Elem, err = NewTypeType(exprType.Elt)
		return
	case *ast.MapType:
		// 字典
		typeType.Kind = Map
		var keyType *TypeType
		keyType, err = NewTypeType(exprType.Key)
		if err != nil {
			return
		}
		typeType.Name = keyType.GetDecl()
		typeType.Elem, err = NewTypeType(exprType.Value)
		return
	case *ast.ChanType:
		// 通道
		typeType.Kind = Chan
		typeType.Elem, err = NewTypeType(exprType.Value)
		return
	case *ast.InterfaceType:
		typeType.Kind = Interface
		typeType.Name = "interface{...}"
		return
	case *ast.IndexExpr:
		// 附带有1个泛型参数的S[X1 Y1]或者S[X1]结构的描述，目前是用于描述泛型的`类型形参`或者`类型实参`
		// X才是实际的S[X1 Y1]中S的部分
		newType, err := NewTypeType(exprType.X)
		if err != nil {
			return nil, err
		}
		typeParam, err := NewTypeType(exprType.Index)
		if err != nil {
			return nil, err
		}
		// log.Println(gjson.MustEncodeString(typeParam))
		// 但是附带的TypeParams还是在当前层级
		// 所以这里对两个层级做一个合并，便于直观的理解
		newType.TypeParams = []*TypeType{typeParam}
		return newType, nil
	case *ast.IndexListExpr:
		// 附带有多个泛型参数的S[X1 Y1,X2 Y2]或者S[X1,X2]结构的描述，目前是用于描述泛型的`类型形参`或者`类型实参`
		typeParams := make([]*TypeType, len(exprType.Indices))
		// 解析类型参数的部分
		for i, indice := range exprType.Indices {
			typeParams[i], err = NewTypeType(indice)
			if err != nil {
				return
			}
		}
		// X才是实际的S[X Y]中S的部分
		newType, err := NewTypeType(exprType.X)
		if err != nil {
			return nil, err
		}
		// 但是附带的TypeParams还是在当前层级
		// 所以这里对两个层级做一个合并，便于直观的理解
		newType.TypeParams = typeParams
		// log.Println(gjson.MustEncodeString(newType))
		return newType, nil
	case *ast.ParenExpr:
		return NewTypeType(exprType.X)
	default:
		err = fmt.Errorf("NewTypeType()未处理的astExpr.(type)=%T: %+v", exprType, exprType)
		return
	}
}

// 如果当前类型是声明者则返回的是S[X1 Y1]的结构
// 如果当前类型是接受者则返回的是S[X1]的结构
func (this *TypeType) getTypeParamsString() string {
	if len(this.TypeParams) == 0 {
		return ""
	}
	sb := &strings.Builder{}
	sb.WriteString("[")
	sb.WriteString(this.TypeParams[0].GetDecl())
	for _, typeParam := range this.TypeParams[1:] {
		sb.WriteString(",")
		sb.WriteString(typeParam.GetDecl())
	}
	sb.WriteString("]")
	return sb.String()
}

func (this *TypeType) GetDecl() string {
	switch this.Kind {
	case Star:
		return fmt.Sprintf("*%s%s", this.Elem.GetDecl(), this.getTypeParamsString())
	case Map:
		return fmt.Sprintf("map[%s%s]%s", this.Name, this.getTypeParamsString(), this.Elem.GetDecl())
	case Array:
		return fmt.Sprintf("[]%s", this.Elem.GetDecl())
	case Ellipsis:
		return fmt.Sprintf("...%s", this.Elem.GetDecl())
	case Chan:
		return fmt.Sprintf("chan %s", this.Elem.GetDecl())
	case Func:
		return fmt.Sprintf("func(...)(...)")
	default:
		return fmt.Sprintf("%s%s", this.Name, this.getTypeParamsString())
	}
}
