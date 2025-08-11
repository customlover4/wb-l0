package integrational

import (
	"context"
	"first-task/internal/config"
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	order "first-task/internal/entities/Order"
	payment "first-task/internal/entities/Payment"
	"first-task/internal/storage/postgres"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgres(t *testing.T) {
	t.Parallel()
	pgContainer := SetupTestDB(t)
	defer pgContainer.Terminate(context.Background())

	host, err := pgContainer.Host(context.Background())
	require.NoError(t, err)
	port, err := pgContainer.MappedPort(context.Background(), DBMapped)
	require.NoError(t, err)

	str := postgres.NewPostgres(config.PostgresConfig{
		Host:     host,
		Port:     port.Port(),
		User:     DBUser,
		Password: DBPassword,
		DBName:   DBName,
		SSLMode:  false,
	})
	defer str.Shutdown()

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

	t.Run("create and find order", func(t *testing.T) {
		err := str.Add(testOrder)
		require.NoError(t, err)

		fromDB, err := str.Find(testOrder.OrderUID)
		require.NoError(t, err)

		require.Equal(t, testOrder, fromDB)
	})

	t.Run("load initial data", func(t *testing.T) {
		initialData, err := str.GetInitialData(100)
		require.NoError(t, err)
		require.Equal(t, 1, len(initialData))
		require.Equal(t, initialData[0], testOrder)
	})
}
