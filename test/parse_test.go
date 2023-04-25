package aster

import (
	"log"
	"testing"

	aster "github.com/szyhf/go-aster"
)

func TestParseStruct(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	modelDir := "./data"

	pkgsTyp, err := aster.ParseDir(modelDir, nil)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("\n" + pkgsTyp[0].String())
	if len(pkgsTyp) != 1 {
		t.Fatalf("包数量不符合预期：%d", len(pkgsTyp))
	}

	curPkgType := pkgsTyp[0]
	if len(pkgsTyp) != 1 {
		t.Fatalf("引用数量不符合预期：%d", len(curPkgType.Imports))
	}
	if curPkgType.Imports[0].Alias != `eu` {
		t.Fatalf("解析的包别名不符合预期：%s", curPkgType.Imports[0].Name)
	}
	if curPkgType.Imports[0].Name != `github.com/szyhf/go-aster/test/data/enum` {
		t.Fatalf("解析的包不符合预期：%s", curPkgType.Imports[0].Name)
	}
	// t.Log(curPkgType.Funcs[0].Name)

	if len(curPkgType.Structs) != 6 {
		t.Fatalf("结构体数量不符合预期：%d", len(curPkgType.Structs))
	}
	expectStructs := map[string]string{
		"Like":                            "TableName",
		"User":                            "HEHE",
		"APP":                             "HAHA",
		"LikeGeneric[K comparable,V any]": "TableNameGeneric",
		"APPGeneric[V any]":               "HAHAGeneric",
		"UserGeneric[X any,Y any]":        "HEHEGeneric",
	}
	for _, structTyp := range curPkgType.Structs {
		if structTyp.PackageType != curPkgType {
			t.Fatalf("结构体引用的包类型不是预计的包类型")
		}
		if expMethod, ok := expectStructs[structTyp.GetDeclName()]; !ok {
			t.Fatalf("结构体名称不符合预期：%s", structTyp.GetDeclName())
		} else {
			if len(structTyp.Methods) != 1 {
				t.Fatalf("结构体%s方法数量不符合预期：%d", structTyp.GetDeclName(), len(structTyp.Methods))
			}
			if structTyp.Methods[0].Name != expMethod {
				t.Fatalf("结构体%s方法名称不符合预期：%s", structTyp.GetDeclName(), structTyp.Methods[0].Name)
			}
		}
	}

	expectFuncs := map[string]bool{"Hello": true, "World": true}
	for _, funTyp := range curPkgType.Funcs {
		if _, ok := expectFuncs[funTyp.Name]; !ok {
			t.Fatalf("包内%s方法名称不符合预期：%s", funTyp.Name, funTyp.Name)
		}
	}
}
