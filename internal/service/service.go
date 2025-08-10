package service

import (
	"context"
	"encoding/json"
	"errors"
	"first-task/internal/config"
	order "first-task/internal/entities/Order"
	"first-task/internal/storage/postgres"
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
}

type MessageListener interface {
	FindOrder(orderUID string) (*order.Order, error)
	AddOrder(orderUID *order.Order) error
}

func NewOrderReader(cfg config.KafkaOrdersConfig) *Service {
	return &Service{
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

func (or *Service) ListenMessages(str MessageListener) {
	zap.L().Info("start listening kafka messages")

	for {
		msg, err := or.reader.ReadMessage(context.Background())
		if errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			zap.L().Error("on reading new kafka message")
			continue
		}

		err = or.process(msg, str)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
}

func (or *Service) Shutdown() {
	if err := or.reader.Close(); err != nil {
		zap.L().Error("error on closing reader")
	}
}

func (or *Service) commitMSG(msg kafka.Message) {
	err := or.reader.CommitMessages(context.Background(), msg)
	if err != nil {
		zap.L().Error(
			fmt.Sprintf("wrong on commititing msg %s", err.Error()),
		)
	}
}

func (or *Service) process(msg kafka.Message, str MessageListener) error {
	var err error

	orderUID := string(msg.Key)

	_, err = str.FindOrder(orderUID)
	if err == nil && !errors.Is(err, postgres.ErrNotFound) {
		or.commitMSG(msg)
		return ErrAlreadyExists
	}

	jsonValue := msg.Value

	var ord order.Order
	err = json.Unmarshal(jsonValue, &ord)
	if err != nil {
		or.commitMSG(msg)
		return ErrWrongData
	}

	err = str.AddOrder(&ord)
	if err != nil {
		return ErrInternalStorage
	}

	zap.L().Info(
		fmt.Sprintf("new order succesfully added: %s", orderUID),
	)

	or.commitMSG(msg)
	return nil
}
