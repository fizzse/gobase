package consumer

import (
	"context"
	"errors"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/pkg/mq/kafka"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

/*
 * kafka consumer
 * 多worker 模式
 * offset全提交 FIXME
 */

type KafkaCfg struct {
	Broker  []string `yaml:"broker"`
	GroupId string   `yaml:"groupId"`
}

type Scheduler struct {
	BizCtx   *biz.SampleBiz
	logger   *zap.SugaredLogger
	conf     KafkaCfg
	connects map[string]*Worker
}

func NewScheduler(logger *zap.SugaredLogger, conf KafkaCfg, bizCtx *biz.SampleBiz) (srv *Scheduler, err error) {
	srv = &Scheduler{
		logger:   logger,
		BizCtx:   bizCtx,
		conf:     conf,
		connects: make(map[string]*Worker),
	}
	return
}

func (s *Scheduler) Name() string {
	return "kafkaConsumer"
}

func (s *Scheduler) Run(ctx context.Context) (err error) {
	s.Route()
	if len(s.connects) == 0 {
		return
	}

	eg, ctx := errgroup.WithContext(ctx)
	for _, worker := range s.connects {
		worker := worker
		eg.Go(func() error {
			//go func() {
			//	select {
			//	case <-ctx.Done():
			//		worker.Stop()
			//		return
			//	}
			//}()

			err = worker.Run(ctx)
			return err
		})
	}
	err = eg.Wait()
	return err
}

func (s *Scheduler) Stop() {
	s.logger.Info("Scheduler stopping")
	for _, worker := range s.connects {
		worker.Stop()
	}
	return
}

func (s *Scheduler) Sub(topic string, handler kafka.Handler, workerCnt, bufSize int) (err error) {
	if handler == nil {
		return errors.New("handler can not be nil")
	}

	if s.connects == nil {
		s.connects = make(map[string]*Worker)
	}

	conf := WorkerConfig{
		KafkaCfg:    s.conf,
		WorkerCount: workerCnt,
		BufSize:     bufSize,
		Topic:       topic,
	}

	worker, err := NewWorker(s.logger, conf, handler)
	if err != nil {
		return err
	}

	s.connects[topic] = worker
	return
}

type WorkerConfig struct {
	KafkaCfg

	Topic       string `yaml:"topic"`
	BufSize     int    `yaml:"bufSize"`
	WorkerCount int    `yaml:"workerCount"`
	stopped     bool
}

type Worker struct {
	WorkerConfig

	sub      *kafka.Subscriber
	dataChan chan kafka.Event
	logger   *zap.SugaredLogger
	handler  kafka.Handler
}

func NewWorker(logger *zap.SugaredLogger, conf WorkerConfig, handleFunc kafka.Handler) (*Worker, error) {
	sub := kafka.NewSubscriber(conf.Topic, conf.Broker, kafka.ConsumerGroup(conf.GroupId))
	worker := &Worker{
		WorkerConfig: conf,
		logger:       logger,
		handler:      handleFunc,
		sub:          sub,
		dataChan:     make(chan kafka.Event, conf.BufSize),
	}

	return worker, nil
}

func (w *Worker) RecvMsg(ctx context.Context, msg kafka.Event) error {
	w.dataChan <- msg
	return nil
}

// Run FIXME 日志修改
func (w *Worker) Run(ctx context.Context) (err error) {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error { // 将kafka数据桥接到 chan
		return w.sub.Subscribe(ctx, w.RecvMsg)
	})

	for i := 0; i < w.WorkerCount; i++ {
		i := i

		eg.Go(func() error {
			w.logger.Infow("mq worker run", "name", w.Topic, "id", i)
			for {
				select {
				case <-ctx.Done():
					w.logger.Infow("mq worker exiting", "name", w.Topic, "id", i)
					return nil

				case event := <-w.dataChan:
					_ = w.handler(ctx, event)
				}
			}
		})
	}

	return eg.Wait()
}

func (w *Worker) Stop() {
	if w.stopped {
		return
	}

	w.logger.Infow("mq worker stopping", "name", w.Topic)
	w.stopped = true
	_ = w.sub.Close()
	close(w.dataChan)
}
