package dao

import (
	"context"

	"github.com/fizzse/gobase/internal/gobase/model"
	"github.com/fizzse/gobase/pkg/cache/redis"
	"github.com/fizzse/gobase/pkg/db"
)

type Dao interface {
	Close()
	CreateUser(ctx context.Context, user *model.User) error
	QueryUser(ctx context.Context, cond *model.User) (*model.User, error)
}

type SampleDao struct {
	dbConn    *db.DbCtx
	redisConn *redis.Client
}

func New(dbConn *db.DbCtx, redisConn *redis.Client) (Dao, func(), error) {
	daoCtx := &SampleDao{
		dbConn:    dbConn,
		redisConn: redisConn,
	}

	return daoCtx, daoCtx.Close, nil
}

func (d *SampleDao) Close() {
	//d.dbConn.Close()
	d.redisConn.Close()
}
