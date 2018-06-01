package aster

import (
	"fmt"
	"go/ast"
)

type TypeType struct {
	Name string `json:",omitempty"`
	Kind Kind   `json:",omitempty"`

	Elem *TypeType `json:",omitempty"`
}

func NewTypeType(astExpr ast.Expr) (typeType *TypeType, err error) {
	typeType = &TypeType{}
	for {
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
			typeType.Name = keyType.String()
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
		default:
			err = fmt.Errorf("NewTypeType()未处理的astExpr.(type)=%T", exprType)
			return
		}
	}
}

func (this *TypeType) String() string {
	switch this.Kind {
	case Star:
		return fmt.Sprintf("*%s", this.Elem.String())
	case Map:
		return fmt.Sprintf("map[%s]%s", this.Name, this.Elem.String())
	case Array:
		return fmt.Sprintf("[]%s", this.Elem.String())
	case Ellipsis:
		return fmt.Sprintf("...%s", this.Elem.String())
	case Chan:
		return fmt.Sprintf("chan %s", this.Elem.String())
	case Func:
		return fmt.Sprintf("func(...)(...)")
	default:
		return fmt.Sprintf("%s", this.Name)
	}
}
