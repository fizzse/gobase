package dao

import (
	"context"

	"github.com/fizzse/gobase/internal/gobase/model"
)

func (d *SampleDao) CreateUser(ctx context.Context, user *model.User) error {
	err := d.dbConn.Model(user).Create(user).Error
	return err
}

func (d *SampleDao) QueryUser(ctx context.Context, cond *model.User) (*model.User, error) {
	user := &model.User{}
	queryFilter := d.dbConn.Model(cond)

	if cond.ID != 0 {
		queryFilter = queryFilter.Where("id = ?", cond.ID)
	}

	if cond.Account != "" {
		queryFilter = queryFilter.Where("account = ?", cond.Account)
	}

	err := queryFilter.First(user).Error
	return user, err
}
