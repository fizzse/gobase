package dao

import (
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"

	qerror "github.com/fizzse/gobase/pkg/error"
	basev1 "github.com/fizzse/gobase/protoc/gobase/v1"
)

func (d *SampleDao) IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (d *SampleDao) IsDuplicateKey(err error) bool {
	if err == nil {
		return false
	}

	// mysql 主键冲突错误
	if strings.Contains(err.Error(), "Error 1062") {
		return true
	}

	return false
}

func (d *SampleDao) convertNotFoundError(errCode basev1.ERR_CODE, err error) (ne error) {
	if err == nil {
		return err
	}

	if d.IsNotFoundError(err) {
		ne = qerror.NotFound(errCode.String(), err.Error())
	} else {
		ne = qerror.ServiceUnavailable(basev1.ERR_CODE_SERVICE_UNAVAILABLE.String(), err.Error())
	}

	ne = errors.Wrap(ne, err.Error()) // 附加原始错误信息
	return
}

func (d *SampleDao) convertDuplicateKeyError(errCode basev1.ERR_CODE, err error) (ne error) {
	if err == nil {
		return err
	}

	if d.IsDuplicateKey(err) {
		ne = qerror.Conflict(errCode.String(), err.Error())
	} else {
		ne = qerror.ServiceUnavailable(basev1.ERR_CODE_SERVICE_UNAVAILABLE.String(), err.Error())
	}

	ne = errors.Wrap(ne, err.Error()) // 附加原始错误信息
	return
}
