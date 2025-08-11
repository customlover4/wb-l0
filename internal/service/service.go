package service

import (
	"context"
	"encoding/json"
	"errors"
	"first-task/internal/config"
	order "first-task/internal/entities/Order"
	"fmt"
	"io"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var ErrAlreadyExists = errors.New("order already exists")
var ErrWrongData = errors.New("wrong data in kafka message")
var ErrInternalStorage = errors.New("error in storage, can't add order")

type OrderReader interface {
	ReadMessage(context.Context) (kafka.Message, error)
	Close() error
	CommitMessages(context.Context, ...kafka.Message) error
}

type Service struct {
	reader OrderReader
	out    chan *order.Order
}

type MessageListener interface {
	FindOrder(orderUID string) (*order.Order, error)
	AddOrder(orderUID *order.Order) error
}

func NewOrderReader(c chan *order.Order, cfg config.KafkaOrdersConfig) *Service {
	return &Service{
		out: c,

		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers:   cfg.Brokers,
				Topic:     cfg.Topic,
				Partition: 0,
				MinBytes:  cfg.MinBytes,
				MaxBytes:  cfg.MaxBytes,
				GroupID:   cfg.GroupID,
			},
		),
	}
}

func (or *Service) ListenMessages(ctx context.Context) {
	zap.L().Info("start listening kafka messages")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := or.reader.ReadMessage(context.Background())
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				zap.L().Error("on reading new kafka message: " + err.Error())
				continue
			}

			or.process(ctx, msg)
		}
	}
}

func (or *Service) Shutdown() {
	if err := or.reader.Close(); err != nil {
		zap.L().Error("error on closing reader")
	}
	close(or.out)
}

func (or *Service) commitMSG(msg kafka.Message) {
	err := or.reader.CommitMessages(context.Background(), msg)
	if err != nil {
		zap.L().Error(
			fmt.Sprintf("wrong on commititing msg %s", err.Error()),
		)
	}
}

func (or *Service) process(ctx context.Context, msg kafka.Message) error {
	var err error
	jsonValue := msg.Value

	var ord order.Order
	err = json.Unmarshal(jsonValue, &ord)
	if err != nil {
		or.commitMSG(msg)
		return ErrWrongData
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case or.out <- &ord:
			zap.L().Info("new order succesfully get")
			or.commitMSG(msg)
			return nil
		}
	}
}
