package service

import (
	"context"
	"encoding/json"
	"errors"
	"first-task/internal/config"
	order "first-task/internal/entities/Order"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

var ErrWrongData = errors.New("can't unmarshal json from kafka msg")
var ErrNotValidData = errors.New("not valid data")

type OrderReader interface {
	ReadMessage(context.Context) (kafka.Message, error)
	Close() error
	CommitMessages(context.Context, ...kafka.Message) error
}

type Service struct {
	reader   OrderReader
	str      OrderAdder
	validate *validator.Validate
	cfg      config.KafkaOrdersConfig
}

type OrderAdder interface {
	AddOrder(*order.Order) error
}

func NewOrderReader(str OrderAdder, cfg config.KafkaOrdersConfig) *Service {
	return &Service{
		validate: validator.New(),
		str:      str,
		cfg:      cfg,

		reader: newReader(cfg),
	}
}

func (s *Service) ListenMessages(ctx context.Context) {
	zap.L().Info("start listening kafka messages")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := s.reader.ReadMessage(context.Background())
			if err != nil {
				zap.L().Error("kafka down: " + err.Error())
				msg = s.retryKafka(ctx)
				zap.L().Info("kafka is up")
			}

			err = s.process(ctx, msg)
			if err != nil {
				zap.L().Error("on processing new order: " + err.Error())
			} else {
				zap.L().Info("new order succesfully added")
			} // TODO: delete later
		}
	}
}

func (s *Service) process(ctx context.Context, msg kafka.Message) error {
	const op = "internal.service.process"

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			jsonValue := msg.Value
			var ord order.Order
			err := json.Unmarshal(jsonValue, &ord)
			if err != nil {
				s.commitMSG(msg)
				return fmt.Errorf("%s: %w", op, errors.Join(ErrWrongData, err))
			}

			err = s.validate.Struct(ord)
			if err != nil {
				s.commitMSG(msg)
				return fmt.Errorf("%s: %w", op, errors.Join(ErrNotValidData, err))
			}

			err = s.str.AddOrder(&ord)
			if err != nil {
				zap.L().Error("err on adding new order to db" + err.Error())
				s.retryDB(&ord)
			}

			s.commitMSG(msg)
			return nil
		}
	}
}

func newReader(cfg config.KafkaOrdersConfig) *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:  cfg.Brokers,
			Topic:    cfg.Topic,
			MinBytes: cfg.MinBytes,
			MaxBytes: cfg.MaxBytes,
			GroupID:  cfg.GroupID,
		},
	)
}

func (s *Service) commitMSG(msg kafka.Message) {
	err := s.reader.CommitMessages(context.Background(), msg)
	if err != nil {
		for {
			for i := 0; i < 5; i++ {
				// s.reader.Close()

				// r := newReader(s.cfg)

				time.Sleep(time.Second * 10)

				// s.reader = r

				err := s.reader.CommitMessages(context.Background(), msg)
				if err == nil {
					return
				}

				zap.L().Error("can't commit message(kafka): " + err.Error())
			}
		}
	}
}

func (s *Service) retryDB(ord *order.Order) {
	for {
		for i := 0; i < 5; i++ {
			err := s.str.AddOrder(ord)
			if err == nil {
				zap.L().Info("DB retrying success")
				return
			}
			time.Sleep(time.Second * 10)
			zap.L().Error(
				"DB still down, retrying again...\n" + err.Error(),
			)
		}
		zap.L().Error("So much attemps retry DB. Waiting 5 minutes and try again.")
		time.Sleep(time.Minute * 5)
	}
}

func (s *Service) retryKafka(ctx context.Context) kafka.Message {
	for {
		select {
		case <-ctx.Done():
			return kafka.Message{}
		default:
			// s.reader.Close()
			// s.reader := newReader(s.cfg)

			time.Sleep(time.Second * 10)

			msg, err := s.reader.ReadMessage(context.Background())
			if err == nil {
				return msg
			}
			zap.L().Error("kafka still down, retry again... | Err: " + err.Error())
		}
	}
}

func (s *Service) Shutdown() {
	if err := s.reader.Close(); err != nil {
		zap.L().Error("error on closing reader")
	}
}
