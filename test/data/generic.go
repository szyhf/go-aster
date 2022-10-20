package data

import (
	eu "github.com/szyhf/go-aster/test/data/enum"
)

// 这是个测试用的结构体
// 多个类型参数嵌入一个单类型参数
type LikeGeneric[K comparable, V any] struct {
	ID              int64              `json:",omitempty"`
	Status          eu.StatusID        `json:",omitempty"`
	Type            eu.LikeTypeID      `json:",omitempty"`
	RefID           int64              `json:",omitempty"`
	Liker           *User              `json:",omitempty"`
	CreateTimestamp int64              `json:",omitempty"`
	Author          *UserGeneric[K, V] `json:",omitempty"`
	APP             *APPGeneric[K]     `json:",omitempty"`
}

func (LikeGenericRecever *LikeGeneric[K, V]) TableNameGeneric() string {
	return "LikeGeneric"
}

// 省略连续声明的复合参数
// 解析后会处理成[X any, Y any]
type UserGeneric[X, Y any] struct {
	ID     int64       `json:",omitempty"`
	Status eu.StatusID `json:",omitempty"`
	Name   string      `json:",omitempty"`
}

func (u *UserGeneric[X, Y]) HEHEGeneric() string {
	return "userGeneric"
}

// 单个类型参数
type APPGeneric[V any] struct {
	ID     int64       `json:",omitempty"`
	Status eu.StatusID `json:",omitempty"`
	Name   string      `json:",omitempty"`
}

func (a *APPGeneric[V]) HAHAGeneric() string {
	return "appGeneric"
}
