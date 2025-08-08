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
	"fmt"

	"go.uber.org/zap"
)

func main() {
	configFile := flag.String("c", "./config.yml", ".yml config file")

	flag.Parse()

	cfg := config.MustLoad(*configFile)
	fmt.Println(cfg)
	str := storage.NewStorage(
		redisStorage.NewRedisStorage(cfg.RedisConfig),
		postgres.NewPostgres(cfg.PostgresConfig),
	)
	zap.ReplaceGlobals(logger.SetupLogger())

	go str.LoadInitialData(cfg.InitialDataSize)
	go webapp.StartWeb(str, cfg.WebConfig)

	service.ListenMessages(str, cfg.KafkaOrdersConfig)
}
