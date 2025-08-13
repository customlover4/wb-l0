package main

import (
	"context"
	"encoding/json"
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	order "first-task/internal/entities/Order"
	payment "first-task/internal/entities/Payment"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
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
			Phone:   "+79000000000",
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
	topic := "orders_new_event"

	brokers := []string{"localhost:9092", "localhost:9094", "localhost:9096", "localhost:9098"}

	for _, b := range brokers {
		conn, err := kafka.Dial("tcp", b)
		if err != nil {
			panic(err.Error() + " on connection")
		}
		defer conn.Close()

		controller, err := conn.Controller()
		if err != nil {
			panic(err.Error())
		}

		controllerConn, err := kafka.Dial(
			"tcp",
			net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)),
		)
		if err != nil {
			panic(err.Error())
		}
		defer controllerConn.Close()

		topicConfig := kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     6,
			ReplicationFactor: 3,
			ConfigEntries: []kafka.ConfigEntry{
				{
					ConfigName:  "min.insync.replicas",
					ConfigValue: "2",
				},
			},
		}
		err = controllerConn.CreateTopics(topicConfig)
		if err != nil {
			panic(err.Error() + " on creating topic")
		}
	}

	time.Sleep(time.Second * 3)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:      brokers,
		Topic:        topic,
		MaxAttempts:  10,
		BatchTimeout: 100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
	})
	defer writer.Close()

	for range time.NewTicker(time.Millisecond).C {
		ord := NewTESTOrder(GenerateRandomString(30))

		kafkaValue, err := json.MarshalIndent(ord, " ", "  ")
		if err != nil {
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
			fmt.Println(err.Error())
			continue
		}

		fmt.Println("send new message")
	}
}
