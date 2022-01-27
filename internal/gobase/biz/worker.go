package biz

import (
	"context"

	"github.com/fizzse/gobase/pkg/mq/kafka"
)

func (b *SampleBiz) DealMsg(ctx context.Context, msg kafka.Event) (err error) { // 回调函数
	return
}
