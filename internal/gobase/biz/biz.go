package biz

import "github.com/fizzse/gobase/internal/gobase/dao"

//
type IBiz interface {
	CreateUser(user *User) (*User, error)
}

func NewIBiz() IBiz {
	return &SampleBiz{}
}

type SampleBiz struct {
	daoCtx dao.IDao
}
