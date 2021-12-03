package dao

import (
	"context"

	"github.com/fizzse/gobase/internal/gobase/model"
	pbBasev1 "github.com/fizzse/gobase/protoc/gobase/v1"
)

func (d *SampleDao) CreateUser(ctx context.Context, user *model.User) (err error) {
	err = d.dbConn.GetConn().Model(user).Create(user).Error
	return
}

func (d *SampleDao) QueryUser(ctx context.Context, cond *model.User) (user *model.User, err error) {
	queryFilter := d.dbConn.GetConn().Model(cond)

	if cond.ID != 0 {
		queryFilter = queryFilter.Where("id = ?", cond.ID)
	}

	if cond.Account != "" {
		queryFilter = queryFilter.Where("account = ?", cond.Account)
	}

	err = queryFilter.First(user).Error
	err = d.convertNotFoundError(pbBasev1.ERR_CODE_USER_NOT_EXISTS, err)
	return user, err
}
