package main

import (
	"context"
	"encoding/json"
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	order "first-task/internal/entities/Order"
	payment "first-task/internal/entities/Payment"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := ""

	for i := 0; i < length; i++ {
		result += string(charset[rand.Intn(len(charset))])
	}

	return result
}

func NewTESTOrder(orderUID string) *order.Order {
	return &order.Order{
		OrderUID:    orderUID,
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: delivery.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: payment.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []item.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NMID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "45",
				TotalPrice:  317,
				NMID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
			{
				ChrtID:      1,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "42",
				TotalPrice:  317,
				NMID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SMID:              99,
		DateCreated:       "2021-11-26T06:22:19Z",
		OOFShard:          "1",
	}
}

func main() {
	conn, err := kafka.Dial("tcp", "localhost:9092")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	topic := "orders"
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	}
	err = conn.CreateTopics(topicConfig)
	if err != nil {
		panic(err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
	})
	defer writer.Close()

	for range time.NewTicker(time.Second * 20).C {
		ord := NewTESTOrder(GenerateRandomString(30))

		kafkaValue, err := json.MarshalIndent(ord, " ", "  ")
		if err != nil {
			zap.L().Error("on marshaling new test order for kafka")
			continue
		}

		err = writer.WriteMessages(
			context.Background(),
			kafka.Message{
				Key:   []byte(ord.OrderUID),
				Value: kafkaValue,
			},
		)

		if err != nil {
			zap.L().Error("Error on writing new order to kafka")
		}
	}
}
