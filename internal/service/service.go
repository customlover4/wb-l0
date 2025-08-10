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

type OrderReader struct {
	reader *kafka.Reader
}

type MessageListener interface {
	FindOrder(orderUID string) (*order.Order, error)
	AddOrder(orderUID *order.Order) error
}

func NewOrderReader(cfg config.KafkaOrdersConfig) *OrderReader {
	return &OrderReader{
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

func (or *OrderReader) commitMSG(reader *kafka.Reader, msg kafka.Message) {
	err := reader.CommitMessages(context.Background(), msg)
	if err != nil {
		zap.L().Error(
			fmt.Sprintf("wrong on commititing msg %s", err.Error()),
		)
	}
}

func (or *OrderReader) ListenMessages(str MessageListener) {
	zap.L().Info("start listening kafka messages")

	for {
		msg, err := or.reader.ReadMessage(context.Background())
		if errors.Is(err, io.EOF) {
			return
		} else if err != nil {
			zap.L().Error("on reading new kafka message")
			continue
		}

		orderUID := string(msg.Key)

		_, err = str.FindOrder(orderUID)
		if err == nil && !errors.Is(err, postgres.ErrNotFound) {
			zap.L().Info("Order allready exists and try add to db again")
			or.commitMSG(or.reader, msg)
			continue
		}

		jsonValue := msg.Value

		var ord order.Order
		err = json.Unmarshal(jsonValue, &ord)
		if err != nil {
			zap.L().Error("wrong json")
			or.commitMSG(or.reader, msg)
			continue
		}

		err = str.AddOrder(&ord)
		if err != nil {
			zap.L().Error(
				fmt.Sprintf(
					"error on adding to db: %s", err.Error(),
				),
			)
			continue
		}

		zap.L().Info(
			fmt.Sprintf("new order succesfully added: %s", orderUID),
		)

		or.commitMSG(or.reader, msg)
	}
}

func (or *OrderReader) Shutdown() {
	if err := or.reader.Close(); err != nil {
		zap.L().Error("error on closing reader")
	}
}
