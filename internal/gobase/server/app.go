package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type App struct {
	RestServer *http.Server
}

func NewApp(h *http.Server) (app *App, closeFunc func(), err error) {
	app = &App{
		RestServer: h,
	}

	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		if err := h.Shutdown(ctx); err != nil {
			log.Fatalf("httpSrv.Shutdown error(%v)", err)
		}
		cancel()
	}
	return
}
