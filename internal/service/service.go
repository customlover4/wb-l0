package kafkalistener

import (
	"context"
	"encoding/json"
	"errors"
	"first-task/internal/config"
	order "first-task/internal/entities/Order"
	"first-task/internal/storage"
	"first-task/internal/storage/postgres"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	Topic = "orders"
)

func CommitMSG(reader *kafka.Reader, msg kafka.Message) {
	err := reader.CommitMessages(context.Background(), msg)
	if err != nil {
		zap.L().Error(
			fmt.Sprintf("wrong on commititing msg %s", err.Error()),
		)
	}
}

func ListenMessages(str storage.Storager, cfg config.KafkaOrdersConfig) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   cfg.Brokers,
		Topic:     cfg.Topic,
		Partition: 0,
		MinBytes:  cfg.MinBytes,
		MaxBytes:  cfg.MaxBytes,
		GroupID:   cfg.GroupID,
	})
	defer reader.Close()

	zap.L().Info("start listening kafka messages")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			zap.L().Error("on reading new kafka message")
			continue
		}

		orderUID := string(msg.Key)

		_, err = str.FindOrder(orderUID)
		if err == nil && !errors.Is(err, postgres.ErrNotFound) {
			zap.L().Info("Order allready exists and try add to db again")
			CommitMSG(reader, msg)
			continue
		}

		jsonValue := msg.Value

		var ord order.Order
		err = json.Unmarshal(jsonValue, &ord)
		if err != nil {
			zap.L().Error(
				fmt.Sprintf("wrong json in %s topic(kafka)", Topic),
			)
			CommitMSG(reader, msg)
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

		CommitMSG(reader, msg)
	}
}
