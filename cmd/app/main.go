package main

import (
	"first-task/internal/client"
	"first-task/internal/config"
	"first-task/pkg/logger"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "first-task/docs"

	"go.uber.org/zap"
)

// @title Orders API
// @version 1.0
// @description This is simple service for parsing orders
// @termsOfService http://swagger.io/terms/

// @host localhost:8080
// @BasePath /

func main() {
	configFile := flag.String("c", "./config.yml", ".yml config file")
	flag.Parse()

	zap.ReplaceGlobals(logger.SetupLogger())

	cfg := config.MustLoad(*configFile)

	c := client.NewClient(cfg)

	c.Init()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	zap.L().Info("stopping app")
	c.Shutdown()
	time.Sleep(time.Second * 1)
}
