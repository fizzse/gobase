package dao

import "github.com/fizzse/gobase/internal/gobase/model"

func (d *SampleDao) CreateUser(user *model.User) error {
	err := d.dbConn.GetConn().Model(user).Create(user).Error
	return err
}

func (d *SampleDao) QueryUser(cond *model.User) (*model.User, error) {
	user := &model.User{}
	queryFilter := d.dbConn.GetConn().Model(cond)

	if cond.ID != 0 {
		queryFilter = queryFilter.Where("id = ?", cond.ID)
	}

	if cond.Account != "" {
		queryFilter = queryFilter.Where("account = ?", cond.Account)
	}

	err := queryFilter.First(user).Error
	return user, err
}
