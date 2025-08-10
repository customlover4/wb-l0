package service

import (
	"context"
	"errors"
	order "first-task/internal/entities/Order"
	"first-task/internal/storage"
	"testing"

	"github.com/segmentio/kafka-go"
)

type StorageMock struct{}

func (sm *StorageMock) FindOrder(orderUID string) (*order.Order, error) {
	switch orderUID {
	case "testCase1":
		return &order.Order{}, nil // mock getting exists order from database
	case "testCase2":
		return nil, storage.ErrNotFound
	case "testCase3":
		return nil, storage.ErrNotFound
	case "testCase4":
		return nil, storage.ErrNotFound
	}

	return nil, errors.New("write switch-case for every testcase, pls")
}

func (sm *StorageMock) AddOrder(ord *order.Order) error {
	if ord.OrderUID == "testError" {
		return errors.New("mock error")
	}
	return nil
}

type KafkaReaderMock struct{}

func (frm *KafkaReaderMock) CommitMessages(context.Context, ...kafka.Message) error {
	return nil
}

func (frm *KafkaReaderMock) Close() error {
	return nil
}

func (frm *KafkaReaderMock) ReadMessage(context.Context) (kafka.Message, error) {
	return kafka.Message{}, nil
}

type TestCase struct {
	Number int
	Err    error
	msg    kafka.Message
}

var Messages = []TestCase{
	{
		Number: 1,
		Err:    ErrAlreadyExists,
		msg: kafka.Message{
			Key:   []byte("testCase1"),
			Value: []byte("test"),
		},
	},
	{
		Number: 2,
		Err:    ErrWrongData,
		msg: kafka.Message{
			Key:   []byte("testCase2"),
			Value: []byte("test"),
		},
	},
	{
		Number: 3,
		Err:    ErrInternalStorage,
		msg: kafka.Message{
			Key:   []byte("testCase3"),
			Value: testErrorJSON,
		},
	},
	{
		Number: 4,
		Err:    nil,
		msg: kafka.Message{
			Key:   []byte("testCase4"),
			Value: testJSON,
		},
	},
}

func TestMSGProcess(t *testing.T) {
	srv := Service{
		reader: &KafkaReaderMock{},
	}

	for _, v := range Messages {
		err := srv.process(v.msg, &StorageMock{})
		if !errors.Is(err, v.Err) {
			t.Errorf(
				"wrong processing\nwait:%s\nget:%s",
				v.Err.Error(), err.Error(),
			)
			continue
		} else {
			t.Logf("Success, Testcase: %d", v.Number)
		}
	}
}

var testErrorJSON = []byte(`{
   "order_uid": "testError",
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

var testJSON = []byte(`{
   "order_uid": "test",
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
