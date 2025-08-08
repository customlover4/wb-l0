package webapp

import (
	"first-task/internal/config"
	"first-task/internal/storage"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type App struct {
	storage.Storage
}

func StartWeb(str storage.Storager, cw config.WebConfig) {
	server := http.Server{
		Addr:         fmt.Sprintf("%s:%s", cw.Host, cw.Port),
		Handler:      Handle(str),
		ReadTimeout:  cw.ReadTimeout,
		WriteTimeout: cw.WriteTimeout,
	}
	zap.L().Info("Starting web site on localhost:8080...")
	if err := server.ListenAndServe(); err != nil {
		zap.L().Fatal("Error on listening and serve(http): " + err.Error())
	}
}
