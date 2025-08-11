package integrational

import (
	"context"
	"first-task/internal/config"
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	order "first-task/internal/entities/Order"
	payment "first-task/internal/entities/Payment"
	"first-task/internal/storage/redisStorage"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRedis(t *testing.T) {
	t.Parallel()
	redisContainer := SetupTestRedis(t)
	defer redisContainer.Terminate(context.Background())

	host, err := redisContainer.Host(context.Background())
	require.NoError(t, err)
	port, err := redisContainer.MappedPort(context.Background(), RedisMapped)
	require.NoError(t, err)

	localStorage := redisStorage.NewRedisStorage(config.RedisConfig{
		Host: host,
		Port: port.Port(),
	})
	defer localStorage.Shutdown()

	testOrder := &order.Order{
		OrderUID:    "test",
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

	t.Run("add, find and delete test", func(t *testing.T) {
		localStorage.Add(testOrder)
		dataOrder := localStorage.Find(testOrder.OrderUID)
		require.Equal(t, testOrder, dataOrder)
		localStorage.Delete(testOrder.OrderUID)
		dataOrder = localStorage.Find(testOrder.OrderUID)
		if dataOrder != nil {
			t.Error("order didn't deleted from cache")
		}
	})
}
