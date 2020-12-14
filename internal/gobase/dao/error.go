package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func (d *SampleDao) dealError(err error) error {
	if err == nil {
		return err
	}

	if gorm.IsRecordNotFoundError(err) {
		return err
	}

	return errors.Wrap(err, "")
}
