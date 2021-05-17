package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opentracing/opentracing-go"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/internal/gobase/server/consumer"
	"github.com/fizzse/gobase/internal/gobase/server/rpc"
	"github.com/pkg/errors"
)

// 日志

func NewApp(h *http.Server, g *rpc.Server, worker *consumer.Worker, logger *zap.SugaredLogger, tracer opentracing.Tracer) (app *App, closeFunc func(), err error) {
	_ = tracer

	app = &App{
		logger:         logger,
		RestServer:     h,
		GrpcServer:     g,
		ConsumerWorker: worker,
		Signal:         make(chan os.Signal),
	}

	signal.Notify(app.Signal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	closeFunc = func() {
		close(app.Signal)
	}
	return
}

type App struct {
	logger *zap.SugaredLogger
	bizCtx biz.SampleBiz

	RestServer     *http.Server // http server
	GrpcServer     *rpc.Server
	ConsumerWorker *consumer.Worker //
	Signal         chan os.Signal   // 监听信号 TODO grpc server
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	wg, ctx := errgroup.WithContext(ctx)

	// http
	wg.Go(func() error {
		defer func() {
			a.logger.Infow("rest server safe exiting")
		}()

		go func() {
			select {
			case <-ctx.Done():
				timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer timeoutCancel()
				_ = a.RestServer.Shutdown(timeoutCtx)
				return
			}
		}()

		err := a.RestServer.ListenAndServe()
		err = errors.Wrap(err, "rest error")
		return err
	})

	// grpc server
	wg.Go(func() error {
		go func() {
			select {
			case <-ctx.Done():
				a.GrpcServer.Stop()
			}
		}()

		err := a.GrpcServer.Run()
		err = errors.Wrap(err, "grpc server error")
		return err
	})

	// consumer worker
	wg.Go(func() error {
		go func() {
			select {
			case <-ctx.Done():
				a.ConsumerWorker.Close()
			}
		}()

		err := a.ConsumerWorker.Run(ctx, a.ConsumerWorker.BizCtx.DealMsg)
		err = errors.Wrap(err, "mq data worker error")
		return err
	})

	// signal
	wg.Go(func() error {
		defer func() {
			a.logger.Infow("signal watcher safe exiting")
		}()

		select {
		case <-ctx.Done():
			return nil
		case sign := <-a.Signal:
			cancel()
			return fmt.Errorf("recv signal: %v", sign)
		}
	})

	a.logger.Info("app run")
	return wg.Wait()
}
