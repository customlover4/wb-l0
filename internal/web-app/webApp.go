package webapp

import (
	"context"
	"errors"
	"first-task/internal/config"
	"first-task/internal/web-app/handlers"
	"fmt"
	"net/http"

	httpSwager "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type WebApp struct {
	server *http.Server
}

func NewWebApp() *WebApp {
	return &WebApp{}
}

func (wa *WebApp) CreateServer(str handlers.OrderGetter, cw config.WebConfig) {
	mux := http.NewServeMux()

	// swagger
	mux.HandleFunc("/swagger/", httpSwager.WrapHandler)

	mux.HandleFunc("GET /order/{order_uid}", handlers.FindOrderAPI(str))
	mux.HandleFunc("GET /find-order", handlers.FindOrder(str))
	mux.HandleFunc("/", handlers.MainPage())

	wa.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cw.Host, cw.Port),
		Handler:      mux,
		ReadTimeout:  cw.ReadTimeout,
		WriteTimeout: cw.WriteTimeout,
	}
}

func (wa *WebApp) StartServer() {
	err := wa.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		zap.L().Fatal("Error on listening and serve(http): " + err.Error())
	}
}

func (wa *WebApp) Shutdown() {
	if err := wa.server.Shutdown(context.Background()); err != nil {
		zap.L().Error("Error on shutdown server(http): " + err.Error())
	}
}
