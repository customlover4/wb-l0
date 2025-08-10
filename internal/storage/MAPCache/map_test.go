package mapcache

import (
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	order "first-task/internal/entities/Order"
	payment "first-task/internal/entities/Payment"
	"testing"
	"time"
)

func TestAddOrder(t *testing.T) {
	str := NewMapStorage()
	str.Add(testOrder("test"))

	if _, ok := str.Cache[orderUID]; !ok {
		t.Error("value doesn't added")
	}
}

func TestGetOrder(t *testing.T) {
	str := NewMapStorage()
	str.Cache[orderUID] = MAPData{
		Expired: time.Now().Add(time.Hour).Unix(),
		Data:    testOrder("test"),
	}

	ord := str.Find(orderUID)
	if ord == nil {
		t.Error("can't find added order")
	}
}

func TestDeleteOrder(t *testing.T) {
	str := NewMapStorage()
	str.Cache[orderUID] = MAPData{
		Expired: time.Now().Add(time.Hour).Unix(),
		Data:    testOrder("test"),
	}

	str.Delete(orderUID)

	if _, ok := str.Cache[orderUID]; ok {
		t.Error("can't delete value from cache")
	}
}

func TestLoadInitialCache(t *testing.T) {
	str := NewMapStorage()

	v := []*order.Order{
		testOrder("test1"), testOrder("test2"),
		testOrder("test3"), testOrder("test4"),
	}

	str.LoadInitialCache(v)

	if len(v) != len(str.Cache) {
		t.Error("can't load all initial values")
	}
}

func TestClean(t *testing.T) {
	str := &MAPStorage{
		Cache: make(map[string]MAPData),
	}
	str.Cache[orderUID] = MAPData{
		Expired: time.Now().Add(time.Second).Unix(),
		Data:    testOrder("test"),
	}
	time.Sleep(time.Second)
	str.clean()
	if _, ok := str.Cache[orderUID]; ok {
		t.Error("doesn't clean expired values")
	}
}

const orderUID = "test"

func testOrder(orderUID string) *order.Order {
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
