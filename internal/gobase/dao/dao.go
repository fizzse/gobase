package dao

import (
	"github.com/fizzse/gobase/internal/gobase/model"
	"github.com/fizzse/gobase/pkg/db"
)

type IDao interface {
	CreateUser(user *model.User) error
	QueryUser(cond *model.User) (*model.User, error)
}

type SampleDao struct {
	dbConn *db.Ctx
}

func NewIDao() IDao {
	return &SampleDao{}
}
