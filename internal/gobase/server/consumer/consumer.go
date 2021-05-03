package consumer

import (
	"context"

	"github.com/fizzse/gobase/internal/gobase/biz"

	"go.uber.org/zap"

	"github.com/fizzse/gobase/pkg/mq/kafka"
	"golang.org/x/sync/errgroup"
)

/*
 * kafka consumer
 * 多worker 模式
 * offset全提交 FIXME
 */

type WorkerConfig struct {
	Broker      []string `yaml:"broker"`
	Topic       string   `yaml:"topic"`
	GroupId     string   `yaml:"groupId"`
	Name        string   `yaml:"name"`
	BufSize     int      `yaml:"bufSize"`
	WorkerCount int      `yaml:"workerCount"`
}

type Worker struct {
	name        string
	topic       string
	groupId     string
	sub         *kafka.Subscriber
	dataChan    chan kafka.Event
	bufSize     int
	workerCount int
	BizCtx      biz.Biz
	logger      *zap.SugaredLogger
}

func NewWorker(logger *zap.SugaredLogger, conf *WorkerConfig, biz *biz.SampleBiz) (*Worker, error) {
	sub := kafka.NewSubscriber(conf.Topic, conf.Broker, kafka.ConsumerGroup(conf.GroupId))
	worker := &Worker{
		logger: logger,
		BizCtx: biz,

		name:        conf.Name,
		topic:       conf.Topic,
		groupId:     conf.GroupId,
		sub:         sub,
		bufSize:     conf.BufSize,
		dataChan:    make(chan kafka.Event, conf.BufSize),
		workerCount: conf.WorkerCount,
	}

	return worker, nil
}

func (w *Worker) RecvMsg(ctx context.Context, msg kafka.Event) error {
	w.dataChan <- msg
	return nil
}

// Run FIXME 日志修改
func (w *Worker) Run(ctx context.Context, handleFunc kafka.Handler) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error { // 将kafka数据桥接到 chan
		return w.sub.Subscribe(ctx, w.RecvMsg)
	})

	for i := 0; i < w.workerCount; i++ {
		i := i

		eg.Go(func() error {
			w.logger.Infow("mq worker run", "name", w.name, "id", i)
			for {
				select {
				case <-ctx.Done():
					w.logger.Infow("mq worker exiting", "name", w.name, "id", i)
					return nil

				case event := <-w.dataChan:
					_ = handleFunc(ctx, event)
				}
			}
		})
	}

	return eg.Wait()
}

func (w *Worker) Close() {
	_ = w.sub.Close()
	close(w.dataChan)
}
