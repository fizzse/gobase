package dao

import (
	"context"

	"github.com/fizzse/gobase/pkg/loader"

	"github.com/fizzse/gobase/internal/gobase/model"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var Provider = wire.NewSet(New, db.NewConn, redis.NewClient, loader.LoadDbConfig, loader.LoadRedisConfig)

type Dao interface {
	Close()
	CreateUser(ctx context.Context, user *model.User) error
	QueryUser(ctx context.Context, cond *model.User) (*model.User, error)
}

type SampleDao struct {
	dbConn    *gorm.DB
	redisConn *redis.Client
}

func New(dbConn *gorm.DB, redisConn *redis.Client) (Dao, func(), error) {
	daoCtx := &SampleDao{
		dbConn:    dbConn,
		redisConn: redisConn,
	}

	return daoCtx, daoCtx.Close, nil
}

func (d *SampleDao) Close() {
	d.dbConn.Close()
	d.redisConn.Close()
}
