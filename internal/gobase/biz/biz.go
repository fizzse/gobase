package biz

import (
	pbBasev1 "github.com/fizzse/gobase/protoc/v1"
	"go.uber.org/zap"

	"github.com/fizzse/gobase/internal/gobase/dao"
)

// 显示声明 SampleBiz 实现了Biz
//var _ Biz = &SampleBiz{}
//
//// GinBiz http interface
//type GinBiz interface {
//	Ping(ginCtx *gin.Context)
//	CreateUserGin(ginCtx *gin.Context)
//}
//
//type Biz interface {
//	GinBiz              // http server
//	pbBase.GobaseServer // grpc
//
//	Close()
//	CreateUser(ctx context.Context, user *CreateUserReq) (*UserInfo, error)
//
//	DealMsg(ctx context.Context, msg kafka.Event) error
//}

// New 返回抽象的接口
//func New(daoCtx dao.Dao) (Biz, func(), error) {
//	bizCtx := &SampleBiz{
//		daoCtx: daoCtx,
//	}
//
//	return bizCtx, bizCtx.Close, nil
//}

// NewInstance 返回具体实例
func NewInstance(daoCtx *dao.SampleDao, logger *zap.SugaredLogger) (*SampleBiz, func(), error) {
	bizCtx := &SampleBiz{
		daoCtx: daoCtx,
		logger: logger,
	}

	return bizCtx, bizCtx.Close, nil
}

type SampleBiz struct {
	//daoCtx   dao.Dao
	daoCtx *dao.SampleDao
	logger *zap.SugaredLogger

	pbBasev1.UnimplementedGobaseServer // FIXME 默认实现所有api
}

func (b *SampleBiz) Close() {
	b.daoCtx.Close()
}

func (b *SampleBiz) Logger() *zap.SugaredLogger {
	return b.logger
}
