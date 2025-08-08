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
	OrderUID          string            `db:"order_uid" json:"order_uid"`
	TrackNumber       string            `db:"track_number" json:"track_number"`
	Entry             string            `db:"entry" json:"entry"`
	Delivery          delivery.Delivery `db:"delivery" json:"delivery"`
	Payment           payment.Payment   `db:"payment" json:"payment"`
	Items             []item.Item       `db:"items" json:"items"`
	Locale            string            `db:"locale" json:"locale"`
	InternalSignature string            `db:"internal_signature" json:"internal_signature"`
	CustomerID        string            `db:"customer_id" json:"customer_id"`
	DeliveryService   string            `db:"delivery_service" json:"delivery_service"`
	ShardKey          string            `db:"shardkey" json:"shardkey"`
	SMID              int64             `db:"sm_id" json:"sm_id"`
	DateCreated       string            `db:"date_created" json:"date_created"`
	OOFShard          string            `db:"oof_shard" json:"oof_shard"`
}

func (o *Order) GetDataForSQLString(LIDDelivery, LIDPayment int64) []any {
	return []any{
		o.OrderUID, o.TrackNumber, o.Entry,
		LIDDelivery, LIDPayment, o.Locale, o.InternalSignature,
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

type Test struct {
	OrderUID string
	Locale   string
}

func (o *Order) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binaryutils.WriteString(buf, o.OrderUID); err != nil {
		return nil, err
	}
	if err := binaryutils.WriteString(buf, o.TrackNumber); err != nil {
		return nil, err
	}
	if err := binaryutils.WriteString(buf, o.Entry); err != nil {
		return nil, err
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

	if err := binaryutils.WriteString(buf, o.Locale); err != nil {
		return nil, err
	}
	if err := binaryutils.WriteString(buf, o.InternalSignature); err != nil {
		return nil, err
	}
	if err := binaryutils.WriteString(buf, o.CustomerID); err != nil {
		return nil, err
	}
	if err := binaryutils.WriteString(buf, o.DeliveryService); err != nil {
		return nil, err
	}
	if err := binaryutils.WriteString(buf, o.ShardKey); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.LittleEndian, o.SMID); err != nil {
		return nil, fmt.Errorf("SMID: %w", err)
	}

	if err := binaryutils.WriteString(buf, o.DateCreated); err != nil {
		return nil, err
	}
	if err := binaryutils.WriteString(buf, o.OOFShard); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (o *Order) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	var err error

	if o.OrderUID, err = binaryutils.ReadString(r); err != nil {
		return err
	}
	if o.TrackNumber, err = binaryutils.ReadString(r); err != nil {
		return err
	}
	if o.Entry, err = binaryutils.ReadString(r); err != nil {
		return err
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

	if o.Locale, err = binaryutils.ReadString(r); err != nil {
		return err
	}
	if o.InternalSignature, err = binaryutils.ReadString(r); err != nil {
		return err
	}
	if o.CustomerID, err = binaryutils.ReadString(r); err != nil {
		return err
	}
	if o.DeliveryService, err = binaryutils.ReadString(r); err != nil {
		return err
	}
	if o.ShardKey, err = binaryutils.ReadString(r); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &o.SMID); err != nil {
		return fmt.Errorf("smid: %w", err)
	}

	if o.DateCreated, err = binaryutils.ReadString(r); err != nil {
		return err
	}
	if o.OOFShard, err = binaryutils.ReadString(r); err != nil {
		return err
	}

	return nil
}
