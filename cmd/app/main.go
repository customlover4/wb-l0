package main

import (
	"first-task/internal/config"
	service "first-task/internal/service"
	"first-task/internal/storage"
	"first-task/internal/storage/postgres"
	"first-task/internal/storage/redisStorage"
	webapp "first-task/internal/web-app"
	"first-task/pkg/logger"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	configFile := flag.String("c", "./config.yml", ".yml config file")
	flag.Parse()

	cfg := config.MustLoad(*configFile)
	zap.ReplaceGlobals(logger.SetupLogger())

	str := storage.NewStorage(
		redisStorage.NewRedisStorage(cfg.RedisConfig),
		postgres.NewPostgres(cfg.PostgresConfig),
	)
	str.LoadInitialData(cfg.InitialDataSize)

	orderReader := service.NewOrderReader(str, cfg.KafkaOrdersConfig)
	go orderReader.ListenMessages()

	wa := webapp.NewWebApp(str, cfg.WebConfig)
	go wa.StartServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	zap.L().Info("stopping app")

	wa.StopServer()
	zap.L().Info("web app is stopped")

	orderReader.Shutdown()
	zap.L().Info("kafkaReader is stopped")

	str.Shutdown()
	zap.L().Info("storage is stopped")

	time.Sleep(time.Second * 5)
}
