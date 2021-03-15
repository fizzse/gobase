package db

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
)

type ClickHouseCfg Config

// 重命名的意义 在于避免使用wire自动生成代码时产生歧义
type ClickhouseCtx struct {
	ctx *gorm.DB
}

func (c *ClickhouseCtx) GetConn() *gorm.DB {
	return c.ctx
}

func NewClickhouseConn(config *ClickHouseCfg) (*ClickhouseCtx, func(), error) {
	conn, clean, err := newClickhouseConn((*Config)(config))
	return &ClickhouseCtx{ctx: conn}, clean, err
}

func newClickhouseConn(config *Config) (*gorm.DB, func(), error) {
	buf := strings.Builder{}
	buf.WriteString(fmt.Sprintf("tcp://%s:%d", config.Address, config.Port))

	p := dsnParam{}
	if config.User != "" {
		p.Set("username", config.User)
		p.Set("password", config.Password)
	}

	if config.DbName != "" {
		p.Set("database", config.DbName)
	}

	if config.DebugModel {
		p.Set("debug", "true")
	}

	// 从库
	if len(config.AltHosts) > 0 {
		var altBuf strings.Builder
		for i, alt := range config.AltHosts {
			if i != 0 {
				altBuf.WriteString(",")
			}
			altBuf.WriteString(alt)
		}

		p.Set("alt_hosts", altBuf.String())
	}

	buf.WriteString("?")
	buf.WriteString(p.Encode())

	source := buf.String()

	log.Println(source)

	conn, err := gorm.Open(clickhouse.Open(source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDb, err := conn.DB()
	if err != nil {
		return nil, nil, err
	}

	if config.MaxIdleConn != 0 {
		sqlDb.SetMaxIdleConns(config.MaxIdleConn)
	}

	if config.MaxOpenConn != 0 {
		sqlDb.SetMaxOpenConns(config.MaxOpenConn)
	}

	return conn, func() {
		sqlDb.Close()
	}, err
}