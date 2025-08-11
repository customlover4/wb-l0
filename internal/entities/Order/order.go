package order

import (
	"bytes"
	"encoding/binary"
	delivery "first-task/internal/entities/Delivery"
	item "first-task/internal/entities/Item"
	payment "first-task/internal/entities/Payment"
	binaryutils "first-task/pkg/utils/binaryUtils"
	"fmt"
)

type Order struct {
	OrderUID          string            `db:"order_uid" json:"order_uid" validate:"required,alphanum"`
	TrackNumber       string            `db:"track_number" json:"track_number" validate:"required"`
	Entry             string            `db:"entry" json:"entry" validate:"required,alpha,uppercase"`
	Delivery          delivery.Delivery `db:"delivery" json:"delivery" validate:"required"`
	Payment           payment.Payment   `db:"payment" json:"payment" validate:"required"`
	Items             []item.Item       `db:"items" json:"items" validate:"required,min=1,dive"`
	Locale            string            `db:"locale" json:"locale" validate:"required,oneof=ru en fr de"`
	InternalSignature string            `db:"internal_signature" json:"internal_signature" validate:"omitempty"`
	CustomerID        string            `db:"customer_id" json:"customer_id" validate:"required"`
	DeliveryService   string            `db:"delivery_service" json:"delivery_service" validate:"required"`
	ShardKey          string            `db:"shardkey" json:"shardkey" validate:"required,numeric"`
	SMID              int64             `db:"sm_id" json:"sm_id" validate:"required,gt=0"`
	DateCreated       string            `db:"date_created" json:"date_created" validate:"required"`
	OOFShard          string            `db:"oof_shard" json:"oof_shard" validate:"required,numeric"`
}

func (o *Order) GetDataForSQLString(DeliveryID, PaymentID int64) []any {
	return []any{
		o.OrderUID, o.TrackNumber, o.Entry,
		DeliveryID, PaymentID, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.ShardKey, o.SMID,
		o.DateCreated, o.OOFShard,
	}
}

func (o *Order) GetDataForSQLStringOrdersItems(orderID int64) []any {
	res := make([]any, 0, len(o.Items)*2)
	for _, v := range o.Items {
		res = append(res, orderID, v.ChrtID)
	}
	return res
}

func (o *Order) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	strFields := []string{
		o.OrderUID, o.TrackNumber, o.Entry, o.Locale, o.InternalSignature,
		o.CustomerID, o.DeliveryService, o.ShardKey, o.DateCreated, o.OOFShard,
	}
	for _, v := range strFields {
		if err := binaryutils.WriteString(buf, v); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(buf, binary.LittleEndian, o.SMID); err != nil {
		return nil, fmt.Errorf("SMID: %w", err)
	}

	deliveryData, err := o.Delivery.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("delivery: %w", err)
	}
	if err := binaryutils.WriteBytesWithLength(buf, deliveryData); err != nil {
		return nil, fmt.Errorf("delivery data: %w", err)
	}

	paymentData, err := o.Payment.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("payment: %w", err)
	}
	if err := binaryutils.WriteBytesWithLength(buf, paymentData); err != nil {
		return nil, fmt.Errorf("payment data: %w", err)
	}

	if err := binary.Write(buf, binary.LittleEndian, uint32(len(o.Items))); err != nil {
		return nil, fmt.Errorf("items count: %w", err)
	}
	for _, item := range o.Items {
		itemData, err := item.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("item: %w", err)
		}
		if err := binaryutils.WriteBytesWithLength(buf, itemData); err != nil {
			return nil, fmt.Errorf("item data: %w", err)
		}
	}

	return buf.Bytes(), nil
}

func (o *Order) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	var err error

	strFields := []*string{
		&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
		&o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.DateCreated, &o.OOFShard,
	}
	for _, v := range strFields {
		str, err := binaryutils.ReadString(r)
		if err != nil {
			return err
		}
		*v = str
	}

	if err := binary.Read(r, binary.LittleEndian, &o.SMID); err != nil {
		return fmt.Errorf("smid: %w", err)
	}

	deliveryData, err := binaryutils.ReadBytesWithLength(r)
	if err != nil {
		return fmt.Errorf("delivery data: %w", err)
	}
	if err := o.Delivery.UnmarshalBinary(deliveryData); err != nil {
		return fmt.Errorf("delivery: %w", err)
	}

	paymentData, err := binaryutils.ReadBytesWithLength(r)
	if err != nil {
		return fmt.Errorf("payment data: %w", err)
	}
	if err := o.Payment.UnmarshalBinary(paymentData); err != nil {
		return fmt.Errorf("payment: %w", err)
	}

	var itemsCount uint32
	if err := binary.Read(r, binary.LittleEndian, &itemsCount); err != nil {
		return fmt.Errorf("items count: %w", err)
	}
	o.Items = make([]item.Item, itemsCount)
	for i := range o.Items {
		itemData, err := binaryutils.ReadBytesWithLength(r)
		if err != nil {
			return fmt.Errorf("item %d data: %w", i, err)
		}
		if err := o.Items[i].UnmarshalBinary(itemData); err != nil {
			return fmt.Errorf("item %d: %w", i, err)
		}
	}

	return nil
}
