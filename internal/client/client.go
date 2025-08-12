package client

import (
	"context"
	"first-task/internal/config"
	order "first-task/internal/entities/Order"
	"first-task/internal/service"
	"first-task/internal/storage"
	"first-task/internal/storage/postgres"
	"first-task/internal/storage/redisStorage"
	webapp "first-task/internal/web-app"
	"first-task/internal/web-app/handlers"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	str Storager
	wa  WebApper
	srv Servicer

	cfg *config.Config
}

type Storager interface {
	AddOrder(ord *order.Order) error
	FindOrder(orderUID string) (*order.Order, error)
	LoadInitialData(size int) error
	Shutdown()
}

type Servicer interface {
	ListenMessages(context.Context)
	Shutdown()
}

type WebApper interface {
	CreateServer(str handlers.OrderGetter, cw config.WebConfig)
	StartServer()
	Shutdown()
}

func NewClient(cfg *config.Config) *Client {
	str := storage.NewStorage(
		redisStorage.NewRedisStorage(cfg.RedisConfig),
		postgres.NewPostgres(cfg.PostgresConfig),
	)
	srv := service.NewOrderReader(str, cfg.KafkaOrdersConfig)
	wa := webapp.NewWebApp()

	return &Client{
		str: str,
		wa:  wa,
		srv: srv,
		cfg: cfg,
	}
}

func (c *Client) Init() {
	err := c.str.LoadInitialData(c.cfg.InitialDataSize)
	if err != nil {
		zap.L().Warn(
			"can't load initial data, skipping this | Err: " + err.Error(),
		)
	}

	serviceCtx, finishService := context.WithCancel(context.Background())
	defer finishService()
	go c.srv.ListenMessages(serviceCtx)

	c.wa.CreateServer(c.str, c.cfg.WebConfig)
	go c.wa.StartServer()

	// gracefull shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	zap.L().Info("stopping app")
	finishService()
	c.Shutdown()
	time.Sleep(time.Second * 1)
}

func (c *Client) Shutdown() {
	c.wa.Shutdown()
	zap.L().Info("web app is stopped")

	c.srv.Shutdown()
	zap.L().Info("kafkaReader is stopped")

	c.str.Shutdown()
	zap.L().Info("storage is stopped")
}
