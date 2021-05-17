package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlCtx struct {
	ctx *gorm.DB
}

func (c *MysqlCtx) GetConn() *gorm.DB {
	return c.ctx
}

type MysqlCfg Config

func NewMysqlConn(config *MysqlCfg) (*MysqlCtx, func(), error) {
	conn, clean, err := newMysqlConn((*Config)(config))
	return &MysqlCtx{ctx: conn}, clean, err
}

func newMysqlConn(option *Config) (*gorm.DB, func(), error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=true&loc=Local",
		option.User,
		option.Password,
		fmt.Sprintf("%s:%d", option.Address, option.Port),
		option.DbName,
		option.Charset,
	)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	if option.DebugModel {
		conn = conn.Debug()
	}

	sqlDb, err := conn.DB()
	if err != nil {
		return nil, nil, err
	}

	if option.MaxIdleConn != 0 {
		sqlDb.SetMaxIdleConns(option.MaxIdleConn)
	}

	if option.MaxOpenConn != 0 {
		sqlDb.SetMaxOpenConns(option.MaxOpenConn)
	}

	return conn, func() {
		sqlDb.Close()
	}, err
}
