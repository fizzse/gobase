package db

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	MysqlDrive     = "mysql"
	PostgresDevice = "postgres"
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

func NewConn(config *Config) (*gorm.DB, func(), error) {
	var db *gorm.DB
	var err error
	switch config.Drive {
	case MysqlDrive:
		db, err = NewMysqlConn(config)
		if err != nil {
			return nil, nil, err
		}

	case PostgresDevice:
		return nil, nil, errors.New("unknown drive")
	default:
		return nil, nil, errors.New("unknown drive")
	}

	db.SingularTable(true)
	db.LogMode(config.DebugModel)
	return db, func() {
		db.Close()
	}, err
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
