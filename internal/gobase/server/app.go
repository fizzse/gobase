package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/fizzse/gobase/internal/gobase/biz"
	"github.com/fizzse/gobase/internal/gobase/server/consumer"
	"github.com/fizzse/gobase/internal/gobase/server/rest"
	"github.com/fizzse/gobase/internal/gobase/server/rpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
)

// 日志

func NewApp(h *rest.Server, g *rpc.Server, scheduler *consumer.Scheduler, logger *zap.SugaredLogger, tracer opentracing.Tracer) (app *App, closeFunc func(), err error) {
	_ = tracer

	app = &App{
		logger:            logger,
		RestServer:        h,
		GrpcServer:        g,
		ConsumerScheduler: scheduler,
		Signal:            make(chan os.Signal),
	}

	// srv 需要实现 IServer
	app.IServers = append(app.IServers, h)
	app.IServers = append(app.IServers, g)
	//app.IServers = append(app.IServers, scheduler)

	signal.Notify(app.Signal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	closeFunc = func() {
		close(app.Signal)
	}
	return
}

type IServer interface {
	Name() string
	Run(ctx context.Context) error
	Stop()
}

type App struct {
	logger *zap.SugaredLogger
	bizCtx biz.SampleBiz

	RestServer        *rest.Server        // http server
	GrpcServer        *rpc.Server         // grpc server
	ConsumerScheduler *consumer.Scheduler // kafka consumer
	IServers          []IServer           // servers
	Signal            chan os.Signal      // 监听信号
}

func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	wg, ctx := errgroup.WithContext(ctx)

	for _, srv := range a.IServers {
		srv := srv
		wg.Go(func() error {
			a.logger.Infof("server %s running", srv.Name())
			defer func() {
				a.logger.Infof("server %s safe exiting", srv.Name())
			}()

			go func() {
				select {
				case <-ctx.Done():
					srv.Stop()
				}
			}()

			err := srv.Run(ctx)
			err = errors.WithMessagef(err, "server: %s exit", srv.Name())
			return err
		})
	}

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

	a.logger.Info("app running")
	return wg.Wait()
}

/*func (a *App) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	wg, ctx := errgroup.WithContext(ctx)

	// http
	wg.Go(func() error {
		defer func() {
			a.logger.Infof("server %s safe exiting", "rest") // logger
		}()

		go func() {
			select {
			case <-ctx.Done():
				a.RestServer.Stop()
				return
			}
		}()

		err := a.RestServer.Run(ctx)
		err = errors.Wrap(err, "rest error")
		return err
	})

	// grpc server
	wg.Go(func() error {
		defer func() {
			a.logger.Infof("server %s safe exiting", "grpc")
		}()

		go func() {
			select {
			case <-ctx.Done():
				a.GrpcServer.Stop()
			}
		}()

		err := a.GrpcServer.Run(ctx)
		err = errors.Wrap(err, "grpc server error")
		return err
	})

	// consumer worker
	wg.Go(func() error {
		defer func() {
			a.logger.Infof("server %s safe exiting", "consumer")
		}()

		go func() {
			select {
			case <-ctx.Done():
				a.ConsumerScheduler.Stop()
			}
		}()

		err := a.ConsumerScheduler.Run(ctx)
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

	a.logger.Info("app running")
	return wg.Wait()
}
*/
