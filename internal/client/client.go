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
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	str Storager
	wa  WebApper
	srv Servicer

	ordersChannel chan *order.Order

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
	out := make(chan *order.Order)

	str := storage.NewStorage(
		redisStorage.NewRedisStorage(cfg.RedisConfig),
		postgres.NewPostgres(cfg.PostgresConfig),
	)
	srv := service.NewOrderReader(out, cfg.KafkaOrdersConfig)
	wa := webapp.NewWebApp()

	return &Client{
		str: str,
		wa:  wa,
		srv: srv,
		cfg: cfg,

		ordersChannel: out,
	}
}

func (c *Client) Init() {
	err := c.str.LoadInitialData(c.cfg.InitialDataSize)
	if err != nil {
		zap.L().Warn(
			"can't load initial data, skipping this | Err: " + err.Error(),
		)
	}
	
	go c.listenMessages()

	c.wa.CreateServer(c.str, c.cfg.WebConfig)
	go c.wa.StartServer()
}

func (c *Client) listenMessages() {
	ctx, finish := context.WithCancel(context.Background())
	defer finish()

	go c.srv.ListenMessages(ctx)

	for v := range c.ordersChannel {
		err := c.str.AddOrder(v)
		if err != nil {
			zap.L().Warn("DB is down, trying retry: " + err.Error())

			c.retry(v)
		}

		zap.L().Info(fmt.Sprintf("new order added: %s", v.OrderUID))
	}
}

func (c *Client) retry(ord *order.Order) {
	for {
		for i := 0; i < 5; i++ {
			err := c.str.AddOrder(ord)
			if err == nil {
				zap.L().Info("DB retrying success")
				return
			}
			time.Sleep(time.Second * 5)
			zap.L().Error(
				"DB still down, retrying again...\n" + err.Error(),
			)
		}
		zap.L().Error("So much attemps retry DB. Waiting 5 minutes and try again.")
		time.Sleep(time.Minute * 5)
	}
}

func (c *Client) Shutdown() {
	c.wa.Shutdown()
	zap.L().Info("web app is stopped")

	c.srv.Shutdown()
	zap.L().Info("kafkaReader is stopped")

	c.str.Shutdown()
	zap.L().Info("storage is stopped")
}
