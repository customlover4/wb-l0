package client

import (
	"first-task/internal/config"
	order "first-task/internal/entities/Order"
	"first-task/internal/service"
	"first-task/internal/storage"
	"first-task/internal/storage/postgres"
	"first-task/internal/storage/redisStorage"
	webapp "first-task/internal/web-app"
	"first-task/internal/web-app/handlers"

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
	LoadInitialData(size int)
	Shutdown()
}

type Servicer interface {
	ListenMessages(str service.MessageListener)
	Shutdown()
}

type WebApper interface {
	CreateServer(str handlers.OrderGetter, cw config.WebConfig)
	StartServer()
	Shutdown()
}

func NewClient(cfgPath string) *Client {
	cfg := config.MustLoad(cfgPath)

	str := storage.NewStorage(
		redisStorage.NewRedisStorage(cfg.RedisConfig),
		postgres.NewPostgres(cfg.PostgresConfig),
	)
	srv := service.NewOrderReader(cfg.KafkaOrdersConfig)
	wa := webapp.NewWebApp()

	return &Client{
		str: str,
		wa:  wa,
		srv: srv,
		cfg: cfg,
	}
}

func (c *Client) Init() {
	c.str.LoadInitialData(c.cfg.InitialDataSize)

	go c.srv.ListenMessages(c.str)

	c.wa.CreateServer(c.str, c.cfg.WebConfig)
	go c.wa.StartServer()
}

func (c *Client) Shutdown() {
	c.wa.Shutdown()
	zap.L().Info("web app is stopped")

	c.srv.Shutdown()
	zap.L().Info("kafkaReader is stopped")

	c.str.Shutdown()
	zap.L().Info("storage is stopped")
}
