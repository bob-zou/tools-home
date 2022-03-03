package di

import (
	"context"
	"net/http"
	"time"
	"tools-home/internal/service"

	"github.com/sirupsen/logrus"
)

//go:generate wire
type App struct {
	svc  *service.Service
	http *http.Server
}

func NewApp(svc *service.Service, h *http.Server) (app *App, closeFunc func(), err error) {
	app = &App{
		svc:  svc,
		http: h,
	}
	closeFunc = func() {
		ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
		if err = h.Shutdown(ctx); err != nil {
			logrus.Infof("httpSrv.Shutdown error(%v)", err)
		}
		cancel()
	}
	return
}
