package dao

import (
	"github.com/fizzse/gobase/internal/gobase/model"
	"github.com/jinzhu/gorm"
)

type IDao interface {
	CreateUser(user *model.User) error
	QueryUser(cond *model.User) (*model.User, error)
}

type SampleDao struct {
	dbConn *gorm.DB
}

func NewIDao() IDao {
	return &SampleDao{}
}
