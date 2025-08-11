package handlers

import (
	"errors"
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	order "first-task/internal/entities/Order"
	payment "first-task/internal/entities/Payment"
	"first-task/internal/storage"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainPage(t *testing.T) {
	r := http.NewServeMux()
	r.HandleFunc("/", MainPage())

	srv := httptest.NewServer(r)
	defer srv.Close()

	resp, err := http.Get(fmt.Sprintf("%s/", srv.URL))
	if err != nil {
		t.Error(err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf(
			"wrong response \nget: %d\nwait: %d",
			resp.StatusCode, http.StatusOK,
		)
	}
}

type TestCase struct {
	Arg  string
	Code int
}

type StorageMock struct{}

func (sm *StorageMock) FindOrder(orderUID string) (*order.Order, error) {
	switch orderUID {
	case "found":
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
		}, nil
	case "not_found":
		return nil, storage.ErrNotFound
	case "wrong_answer":
		return nil, errors.New("test error from db")
	}
	return nil, nil
}

func TestOrderPage(t *testing.T) {
	tests := []TestCase{
		{
			Arg:  "not_found",
			Code: http.StatusNotFound,
		},
		{
			Arg:  "found",
			Code: http.StatusOK,
		},
		{
			Arg:  "wrong_answer",
			Code: http.StatusInternalServerError,
		},
	}

	r := http.NewServeMux()
	r.HandleFunc("/find-order", FindOrder(&StorageMock{}))

	srv := httptest.NewServer(r)
	defer srv.Close()

	for _, v := range tests {
		resp, err := http.Get(fmt.Sprintf(
			"%s/find-order?order_uid=%s", srv.URL, v.Arg),
		)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		if resp.StatusCode != v.Code {
			t.Errorf(
				"wrong response \nget: %d\nwait: %d\nfrom: %s",
				resp.StatusCode, v.Code, resp.Request.URL.String(),
			)
		}
	}
}

func TestOrderAPI(t *testing.T) {
	tests := []TestCase{
		{
			Arg:  "not_found",
			Code: http.StatusNotFound,
		},
		{
			Arg:  "found",
			Code: http.StatusOK,
		},
		{
			Arg:  "wrong_answer",
			Code: http.StatusInternalServerError,
		},
	}

	r := http.NewServeMux()
	r.HandleFunc("/order/{order_uid}", FindOrderAPI(&StorageMock{}))
	
	srv := httptest.NewServer(r)
	defer srv.Close()

	for _, v := range tests {
		resp, err := http.Get(fmt.Sprintf(
			"%s/order/%s", srv.URL, v.Arg),
		)
		if err != nil {
			t.Error(err.Error())
			continue
		}

		if resp.StatusCode != v.Code {
			t.Errorf(
				"wrong response \nget: %d\nwait: %d\nfrom: %s",
				resp.StatusCode, v.Code, resp.Request.URL.String(),
			)
		}
	}
}

// func TestFindOrder(t *testing.T) {

// }
