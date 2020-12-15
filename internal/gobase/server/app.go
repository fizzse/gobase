package server

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	RestServer *http.Server   // http server
	Signal     chan os.Signal // 监听信号 TODO grpc server
}

func NewApp(h *http.Server) (app *App, closeFunc func(), err error) {
	app = &App{
		RestServer: h,
		Signal:     make(chan os.Signal),
	}

	signal.Notify(app.Signal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	closeFunc = func() {
		//ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		//if err := h.Shutdown(ctx); err != nil {
		//	log.Fatalf("httpSrv.Shutdown error(%v)", err)
		//}
		//cancel()
	}
	return
}
