package integrational

import (
	"context"
	"first-task/internal/config"
	order "first-task/internal/entities/Order"
	"first-task/internal/service"
	"first-task/internal/storage"
	"fmt"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

func TestNewOrder(t *testing.T) {
	t.Parallel()

	kafkaContainer := SetupTestKafka(t)
	defer kafkaContainer.Terminate(context.Background())

	KafkaHost, err := kafkaContainer.Host(context.Background())
	require.NoError(t, err)
	KafkaPort, err := kafkaContainer.MappedPort(context.Background(), "9092")
	require.NoError(t, err)

	postgresContainer, database := SetupTestDB(t)
	defer postgresContainer.Terminate(context.Background())

	redisContainer, local := SetupTestRedis(t)
	defer redisContainer.Terminate(context.Background())

	str := storage.NewStorage(local, database)
	defer str.Shutdown()

	t.Run("new order event/get order", func(t *testing.T) {
		brokerAddr := fmt.Sprintf("%s:%s", KafkaHost, KafkaPort.Port())
		ordersC := make(chan *order.Order)
		kfk := service.NewOrderReader(ordersC, config.KafkaOrdersConfig{
			Brokers: []string{
				brokerAddr,
			},
			Topic:    KafkaTopic,
			MinBytes: 1,
			MaxBytes: 10e6,
			GroupID:  "test-group",
		})

		conn, err := kafka.Dial("tcp", brokerAddr)
		require.NoError(t, err)

		topicConfig := kafka.TopicConfig{
			Topic:             KafkaTopic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		}

		err = conn.CreateTopics(topicConfig)
		conn.Close()
		require.NoError(t, err)

		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{brokerAddr},
			Topic:   KafkaTopic,
		})

		orderUID := "integrationaltestaddandfindorder"
		msg := kafka.Message{
			Key:   []byte(orderUID),
			Value: testJSON,
		}

		err = writer.WriteMessages(context.Background(), msg)
		writer.Close()
		require.NoError(t, err)

		ctx, finish := context.WithCancel(context.Background())
		defer finish()
		go kfk.ListenMessages(ctx)
		ord := <-ordersC
		finish()
		kfk.Shutdown()

		err = str.AddOrder(ord)
		require.NoError(t, err)

		DBOrd, err := str.FindOrder(ord.OrderUID)
		require.NoError(t, err)

		require.Equal(t, ord, DBOrd)
	})
}

var testJSON = []byte(`{
   "order_uid": "b563feb7b2b84b6test",
   "track_number": "WBILMTESTTRACK",
   "entry": "WBIL",
   "delivery": {
      "name": "Test Testov",
      "phone": "+9720000000",
      "zip": "2639809",
      "city": "Kiryat Mozkin",
      "address": "Ploshad Mira 15",
      "region": "Kraiot",
      "email": "test@gmail.com"
   },
   "payment": {
      "transaction": "b563feb7b2b84b6test",
      "request_id": "",
      "currency": "USD",
      "provider": "wbpay",
      "amount": 1817,
      "payment_dt": 1637907727,
      "bank": "alpha",
      "delivery_cost": 1500,
      "goods_total": 317,
      "custom_fee": 0
   },
   "items": [
      {
         "chrt_id": 9934930,
         "track_number": "WBILMTESTTRACK",
         "price": 453,
         "rid": "ab4219087a764ae0btest",
         "name": "Mascaras",
         "sale": 30,
         "size": "0",
         "total_price": 317,
         "nm_id": 2389212,
         "brand": "Vivienne Sabo",
         "status": 202
      }
   ],
   "locale": "en",
   "internal_signature": "",
   "customer_id": "test",
   "delivery_service": "meest",
   "shardkey": "9",
   "sm_id": 99,
   "date_created": "2021-11-26T06:22:19Z",
   "oof_shard": "1"
}`)
