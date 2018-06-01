package data

import (
	eu "github.com/szyhf/go-aster/test/data/enum"
)

// 这是个测试用的结构体
type Like struct {
	ID              int64         `json:",omitempty"`
	Status          eu.StatusID   `json:",omitempty"`
	Type            eu.LikeTypeID `json:",omitempty"`
	RefID           int64         `json:",omitempty"`
	Liker           *User         `json:",omitempty"`
	CreateTimestamp int64         `json:",omitempty"`
	Author          *User         `json:",omitempty"`
	APP             *APP          `json:",omitempty"`
}

func (l *Like) TableName() string {
	return "like"
}

type User struct {
	ID     int64       `json:",omitempty"`
	Status eu.StatusID `json:",omitempty"`
	Name   string      `json:",omitempty"`
}

func (u *User) HEHE() string {
	return "user"
}

type APP struct {
	ID     int64       `json:",omitempty"`
	Status eu.StatusID `json:",omitempty"`
	Name   string      `json:",omitempty"`
}

func (a *APP) HAHA() string {
	return "app"
}
