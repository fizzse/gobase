package biz

import (
	"context"

	pbBasev1 "github.com/fizzse/gobase/protoc/v1"
)

func (b *SampleBiz) CreateUser(ctx context.Context, in *pbBasev1.CreateUserReq) (out *pbBasev1.UserInfo, err error) {
	out = &pbBasev1.UserInfo{}
	return
}
