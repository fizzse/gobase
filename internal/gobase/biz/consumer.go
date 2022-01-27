package biz

import (
	"context"

	"github.com/fizzse/gobase/pkg/mq/kafka"
)

func (b *SampleBiz) DealMsg(ctx context.Context, msg kafka.Event) (err error) { // 回调函数
	b.logger.Infow("DealMsg", "msg", msg)
	return
}

func (b *SampleBiz) DealMsg2(ctx context.Context, msg kafka.Event) (err error) { // 回调函数
	b.logger.Infow("DealMsg2", "msg", msg)
	return
}
