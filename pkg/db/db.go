package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Config struct {
	Drive           string        `json:"drive" yaml:"drive" yml:"drive"`                   // 数据库类型
	Address         string        `json:"address" yaml:"address" yml:"address"`             // 地址
	Port            int64         `json:"port" yaml:"port" yml:"port"`                      // 端口
	User            string        `json:"user" yaml:"user" yml:"user"`                      // 用户
	Password        string        `json:"password" yaml:"password" yml:"password"`          // 密码
	DbName          string        `json:"dbName" yaml:"dbName" yml:"dbName"`                // 数据库
	Charset         string        `json:"charset" yaml:"charset" yml:"charset"`             // 字符集
	DebugModel      bool          `json:"debugModel" yaml:"debugModel" yml:"debugModel"`    // 调试模式
	MaxIdleConn     int           `json:"maxIdleConn" yaml:"maxIdleConn" yml:"maxIdleConn"` // 最大空闲链接
	MaxOpenConn     int           `json:"maxOpenConn" yaml:"maxOpenConn" yml:"maxOpenConn"` // 最大打开链接
	ParseTime       bool          `json:"parseTime" yaml:"parseTime"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" yaml:"connMaxLifetime" yml:"connMaxLifetime"` // 链接最大复用时间
	AltHosts        []string      `json:"altHosts" yaml:"altHosts" yml:"altHosts"`                      // 从库
}

type DbCtx struct {
	ctx *gorm.DB
}

const (
	MysqlDrive      = "mysql"
	PostgresDrive   = "postgres"
	ClickhouseDrive = "clickhouse"
)

func NewConn(option *Config) (*DbCtx, func(), error) {
	var db *gorm.DB
	var err error
	var clean func()
	dbCtx := &DbCtx{}
	switch option.Drive {
	case MysqlDrive:
		db, clean, err = newMysqlConn(option)

	case ClickhouseDrive:
		db, clean, err = newClickhouseConn(option)

	case PostgresDrive:
		return nil, nil, errors.New("unknown drive")
	default:
		return nil, nil, errors.New("unknown drive")
	}

	dbCtx.ctx = db
	return dbCtx, clean, err
}

func (c *DbCtx) GetConn() *gorm.DB {
	return c.ctx
}

func (c *DbCtx) Begin() *gorm.DB {
	return c.ctx.Begin()
}

func (c *DbCtx) RollBack() *gorm.DB {
	return c.ctx.Rollback()
}

func (c *DbCtx) Commit() *gorm.DB {
	return c.ctx.Commit()
}
