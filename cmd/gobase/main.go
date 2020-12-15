package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fizzse/gobase/internal/gobase/server"
	"golang.org/x/sync/errgroup"
)

func main() {
	app, cleanFunc, err := server.InitApp()
	if err != nil {
		log.Fatal("init app failed: ", err)
	}

	defer cleanFunc()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		go func() {
			select {
			case <-ctx.Done():
				timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 1*time.Second)
				defer timeoutCancel()
				app.RestServer.Shutdown(timeoutCtx)
				return
			}
		}()

		return app.RestServer.ListenAndServe()
	})

	wg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case sign := <-app.Signal:
			cancel()
			return fmt.Errorf("recv signal: %v", sign)
		}
	})

	if err := wg.Wait(); err != nil {
		log.Println("recv error: ", err)
	}

	log.Println("server stop")
}
