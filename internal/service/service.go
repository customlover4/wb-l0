package service

import (
	"context"
	"encoding/json"
	"errors"
	"first-task/internal/config"
	order "first-task/internal/entities/Order"
	"fmt"
	"io"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var ErrWrongData = errors.New("can't unmarshal json from kafka msg")

type OrderReader interface {
	ReadMessage(context.Context) (kafka.Message, error)
	Close() error
	CommitMessages(context.Context, ...kafka.Message) error
}

type Service struct {
	reader   OrderReader
	out      chan *order.Order
	validate *validator.Validate
}

func NewOrderReader(c chan *order.Order, cfg config.KafkaOrdersConfig) *Service {
	return &Service{
		out:      c,
		validate: validator.New(),

		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers:  cfg.Brokers,
				Topic:    cfg.Topic,
				MinBytes: cfg.MinBytes,
				MaxBytes: cfg.MaxBytes,
				GroupID:  cfg.GroupID,
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
				zap.L().Error(
					"Kafka is down, trying retry... (error: " + err.Error() + ")",
				)

				msg = or.retry()
			}

			err = or.process(ctx, msg)
			if err != nil {
				zap.L().Error("on processing new order: " + err.Error())
			} else {
				zap.L().Info("new order succesfully get")
			}
		}
	}
}

func (or *Service) process(ctx context.Context, msg kafka.Message) error {
	const op = "internal.service.process"

	jsonValue := msg.Value

	var ord order.Order
	err := json.Unmarshal(jsonValue, &ord)
	if err != nil {
		or.commitMSG(msg)
		return fmt.Errorf("%s: %w", op, ErrWrongData)
	}

	err = or.validate.Struct(ord)
	if err != nil {
		or.commitMSG(msg)
		return fmt.Errorf("%s: order validation failed (%w)", op, err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case or.out <- &ord:
			or.commitMSG(msg)
			return nil
		}
	}
}

func (or *Service) commitMSG(msg kafka.Message) {
	err := or.reader.CommitMessages(context.Background(), msg)
	if err != nil {
		for {
			for i := 0; i < 5; i++ {
				err := or.reader.CommitMessages(context.Background(), msg)
				if err == nil {
					return
				}
				zap.L().Error("can't commit message(kafka): " + err.Error())
				time.Sleep(time.Second * 10)
			}
			time.Sleep(time.Minute * 5)
		}
	}
}

func (or *Service) retry() kafka.Message {
	for {
		for i := 0; i < 5; i++ {
			msg, err := or.reader.ReadMessage(context.Background())
			if errors.Is(err, io.EOF) {
				break
			} else if err == nil {
				return msg
			}
			time.Sleep(time.Second * 5)
			zap.L().Error(
				"Kafka still down, retrying again...\n" + err.Error(),
			)
		}
		zap.L().Error("So much attemps retry DB. Waiting 5 minutes and try again.")
		time.Sleep(time.Minute * 5)
	}
}

func (or *Service) Shutdown() {
	if err := or.reader.Close(); err != nil {
		zap.L().Error("error on closing reader")
	}
	close(or.out)
}
