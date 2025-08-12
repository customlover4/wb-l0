package integrational

import (
	"bytes"
	"context"
	"encoding/json"
	"first-task/internal/client"
	"first-task/internal/config"
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	order "first-task/internal/entities/Order"
	payment "first-task/internal/entities/Payment"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func mockOrderToKafka(t *testing.T, brokerAddr string) string {
	conn, err := kafka.Dial("tcp", brokerAddr)
	require.NoError(t, err)
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             KafkaTopic,
		NumPartitions:     3,
		ReplicationFactor: 1,
	}

	err = conn.CreateTopics(topicConfig)
	require.NoError(t, err)
	t.Log("wait for creating topic")
	time.Sleep(time.Second * 5)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddr},
		Topic:   KafkaTopic,
	})
	defer writer.Close()

	jsonData, err := json.Marshal(testOrder)
	require.NoError(t, err)

	t.Log("wait for sending message to kafka")
	err = writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(testOrder.OrderUID),
			Value: jsonData,
		},
	)
	require.NoError(t, err)

	return testOrder.OrderUID
}

func TestNewOrder(t *testing.T) {
	t.Parallel()

	kafkaContainer := SetupTestKafka(t)
	defer kafkaContainer.Terminate(context.Background())

	KafkaHost, err := kafkaContainer.Host(context.Background())
	require.NoError(t, err)
	KafkaPort, err := kafkaContainer.MappedPort(context.Background(), KafkaMapped)
	require.NoError(t, err)

	postgresContainer := SetupTestDB(t)
	defer postgresContainer.Terminate(context.Background())
	PsqlHost, err := postgresContainer.Host(context.Background())
	require.NoError(t, err)
	PsqlPort, err := postgresContainer.MappedPort(context.Background(), DBMapped)
	require.NoError(t, err)

	redisContainer := SetupTestRedis(t)
	defer redisContainer.Terminate(context.Background())

	RedisHost, err := redisContainer.Host(context.Background())
	require.NoError(t, err)
	RedisPort, err := redisContainer.MappedPort(context.Background(), RedisMapped)
	require.NoError(t, err)

	t.Run("new order event and get order", func(t *testing.T) {
		brokerAddr := fmt.Sprintf("%s:%s", KafkaHost, KafkaPort.Port())

		cli := client.NewClient(
			&config.Config{
				WebConfig: config.WebConfig{
					Host: "localhost",
					Port: "8080",
				},
				PostgresConfig: config.PostgresConfig{
					Host:     PsqlHost,
					Port:     PsqlPort.Port(),
					User:     DBUser,
					Password: DBPassword,
					DBName:   DBName,
					SSLMode:  false,
				},
				RedisConfig: config.RedisConfig{
					Host: RedisHost,
					Port: RedisPort.Port(),
				},
				KafkaOrdersConfig: config.KafkaOrdersConfig{
					Brokers: []string{
						brokerAddr,
					},
					Topic:    KafkaTopic,
					MinBytes: 1,
					MaxBytes: 10e6,
					GroupID:  "my-test-group",
				},
			},
		)
		defer cli.Shutdown()

		go cli.Init()

		orderUID := mockOrderToKafka(t, brokerAddr)
		t.Log("Wait for adding to storage...")
		time.Sleep(time.Second * 3)

		resp, err := http.Get(
			"http://localhost:8080/order/error_404",
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)

		resp, err = http.Get(
			fmt.Sprintf("http://localhost:8080/order/%s", orderUID),
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var dt order.Order
		r := bytes.NewBuffer([]byte{})
		r.ReadFrom(resp.Body)
		json.Unmarshal(r.Bytes(), &dt)

		require.Equal(t, testOrder, dt)
	})
}

var testOrder = order.Order{
	OrderUID:    "b563feb7b2b84b6test",
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
