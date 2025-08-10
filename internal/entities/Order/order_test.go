package order

import (
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	payment "first-task/internal/entities/Payment"
	"testing"
)

func CmpOrder(a, b *Order) bool {
	return a.OrderUID == b.OrderUID &&
		a.TrackNumber == b.TrackNumber &&
		a.Entry == b.Entry &&
		len(a.Items) == len(b.Items) &&
		a.Locale == b.Locale &&
		a.InternalSignature == b.InternalSignature &&
		a.CustomerID == b.CustomerID &&
		a.DeliveryService == b.DeliveryService &&
		a.ShardKey == b.ShardKey &&
		a.SMID == b.SMID &&
		a.DateCreated == b.DateCreated &&
		a.OOFShard == b.OOFShard
}

func TestOrderMarshaling(t *testing.T) {
	o := &Order{
		OrderUID:    "test_order_uid",
		TrackNumber: "test_track",
		Entry:       "test_entry",
		Delivery: delivery.Delivery{
			Name:    "test_name",
			Phone:   "test_phone",
			Zip:     "test_zip",
			City:    "test_city",
			Address: "test_address",
			Region:  "test_region",
			Email:   "test@email.com",
		},
		Payment: payment.Payment{
			Transaction: "test_transaction",
			RequestID:   "test_request",
			Currency:    "USD",
			Provider:    "test_provider",
			Amount:      1000,
		},
		Items: []item.Item{
			{
				ChrtID:      123,
				TrackNumber: "test_item_track",
				Price:       500,
				RID:         "test_rid",
				Name:        "test_item",
				Sale:        10,
				Size:        "M",
				TotalPrice:  450,
				NMID:        456,
				Brand:       "test_brand",
				Status:      1,
			},
		},
		Locale:            "en",
		InternalSignature: "test_sig",
		CustomerID:        "test_customer",
		DeliveryService:   "test_service",
		ShardKey:          "test_shard",
		SMID:              1,
		DateCreated:       "2023-01-01T00:00:00Z",
		OOFShard:          "test_oof",
	}

	tmp, err := o.MarshalBinary()
	if err != nil {
		t.Error("failed on marshaling Order: " + err.Error())
		return
	}

	newO := &Order{}
	err = newO.UnmarshalBinary(tmp)
	if err != nil {
		t.Error("failed on unmarshaling Order: " + err.Error())
		return
	}

	if !CmpOrder(o, newO) {
		t.Errorf("Item unmarshaling failed, fields don't match")
		t.Errorf("Original: %+v", o)
		t.Errorf("Unmarshaled: %+v", newO)
	}
}
