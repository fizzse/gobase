package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fizzse/gobase/internal/gobase/biz"

	"go.uber.org/zap"

	"golang.org/x/sync/errgroup"

	"github.com/fizzse/gobase/internal/gobase/server/consumer"
	"github.com/pkg/errors"
)

// 日志

func NewApp(logger *zap.SugaredLogger, h *http.Server, worker *consumer.Worker) (app *App, closeFunc func(), err error) {
	app = &App{
		RestServer:     h,
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
	logger         *zap.SugaredLogger
	bizCtx         biz.Biz
	RestServer     *http.Server     // http server
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

	// consumer worker
	wg.Go(func() error {
		go func() {
			select {
			case <-ctx.Done():
				a.ConsumerWorker.Close()
			}
		}()

		err := a.ConsumerWorker.Run(ctx, a.bizCtx.DealMsg)
		err = errors.Wrap(err, "data worker error")
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
