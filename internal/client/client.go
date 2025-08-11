package client

import (
	"context"
	"encoding/json"
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

	"github.com/segmentio/kafka-go"
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
	LoadInitialData(size int)
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
	c.str.LoadInitialData(c.cfg.InitialDataSize)

	go c.ListenMessages()

	c.wa.CreateServer(c.str, c.cfg.WebConfig)
	go c.wa.StartServer()
}

func (c *Client) ListenMessages() {
	ctx, finish := context.WithCancel(context.Background())
	defer finish()

	go c.srv.ListenMessages(ctx)

	for v := range c.ordersChannel {
		err := c.str.AddOrder(v)
		if err != nil {

			zap.L().Warn("DB is down, trying retry: " + v.OrderUID)

			retryResult := c.retry(v)
			if !retryResult {
				finish()
				c.backoff(v)
				return
			}
		}

		zap.L().Info(fmt.Sprintf("new order added: %s", v.OrderUID))
	}
}

func (c *Client) retry(ord *order.Order) bool {
	successRetry := false
	for i := 0; i < 5; i++ {
		err := c.str.AddOrder(ord)
		if err == nil {
			zap.L().Info("DB retrying success")
			successRetry = true
			break
		}
		time.Sleep(time.Second * 5)
		zap.L().Info("DB still down, retrying again...")
	}

	return successRetry
}

func (c *Client) backoff(ord *order.Order) {
	c.Shutdown()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: c.cfg.KafkaOrdersConfig.Brokers,
		Topic:   c.cfg.KafkaOrdersConfig.Topic,
	})

	jsonData, err := json.Marshal(ord)
	if err != nil {
		zap.L().Error(fmt.Sprintf(
			"can't backoff order: %s", ord.OrderUID,
		))
		return
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(ord.OrderUID),
		Value: jsonData,
	})
	if err != nil {
		zap.L().Error("can't write backoff messages")
		return
	}

	panic("database is down, can't save orders, data returned to kafka")
}

func (c *Client) Shutdown() {
	c.wa.Shutdown()
	zap.L().Info("web app is stopped")

	c.srv.Shutdown()
	zap.L().Info("kafkaReader is stopped")

	c.str.Shutdown()
	zap.L().Info("storage is stopped")
}
