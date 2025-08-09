package webapp

import (
	"context"
	"errors"
	"first-task/internal/config"
	"first-task/internal/storage"
	"first-task/internal/web-app/handlers"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type WebApp struct {
	server *http.Server
	h      handlers.Handler
}

func NewWebApp(str storage.Storager, cw config.WebConfig) *WebApp {
	return &WebApp{
		server: &http.Server{
			Addr:         fmt.Sprintf("%s:%s", cw.Host, cw.Port),
			Handler:      Handle(str),
			ReadTimeout:  cw.ReadTimeout,
			WriteTimeout: cw.WriteTimeout,
		},
		h: *handlers.NewHandler(str),
	}
}

func (wa *WebApp) StartServer() {
	err := wa.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.L().Fatal("Error on listening and serve(http): " + err.Error())
	}
}

func (wa *WebApp) StopServer() {
	if err := wa.server.Shutdown(context.Background()); err != nil {
		zap.L().Error("Error on shutdown server(http): " + err.Error())
	}
}
