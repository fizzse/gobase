package db

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Config struct {
	Drive      string `json:"drive" yaml:"drive" yml:"drive"`
	Address    string `json:"address" yaml:"address" yml:"address"`
	Port       int64  `json:"port" yaml:"port" yml:"port"`
	User       string `json:"user" yaml:"user" yml:"user"`
	Password   string `json:"password" yaml:"password" yml:"password"`
	DbName     string `json:"dbName" yaml:"dbName" yml:"dbName"`
	Charset    string `json:"charset" yaml:"charset" yml:"charset"`
	DebugModel bool   `json:"debugModel" yaml:"debugModel" yml:"debugModel"`
}

type Ctx struct {
	ctx *gorm.DB
}

const (
	MysqlDrive     = "mysql"
	PostgresDevice = "postgres"
)

func NewConn(config *Config) (*Ctx, error) {
	var db *gorm.DB
	var err error
	dbCtx := &Ctx{}
	switch config.Drive {
	case MysqlDrive:
		db, err = NewMysqlConn(config)
		if err != nil {
			return nil, err
		}

	case PostgresDevice:
		return nil, errors.New("unknown drive")
	default:
		return nil, errors.New("unknown drive")
	}

	dbCtx.ctx = db
	dbCtx.ctx.SingularTable(true)
	dbCtx.ctx.LogMode(config.DebugModel)

	return dbCtx, nil
}

func NewMysqlConn(config *Config) (*gorm.DB, error) {
	return gorm.Open(config.Drive, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s",
		config.User,
		config.Password,
		fmt.Sprintf("%s:%d", config.Address, config.Port),
		config.DbName,
		config.Charset,
	))
}

func (c *Ctx) GetConn() *gorm.DB {
	return c.ctx
}

func (c *Ctx) Begin() *gorm.DB {
	return c.ctx.Begin()
}

func (c *Ctx) RollBack() *gorm.DB {
	return c.ctx.Rollback()
}

func (c *Ctx) Commit() *gorm.DB {
	return c.ctx.Commit()
}
