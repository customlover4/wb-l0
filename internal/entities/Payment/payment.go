package payment

import (
	"bytes"
	"encoding/binary"
	binaryutils "first-task/pkg/utils/binaryUtils"
	"fmt"
	"math"
)

type Payment struct {
	Transaction  string  `db:"transaction" json:"transaction" validate:"required"`
	RequestID    string  `db:"request_id" json:"request_id" validate:"omitempty,alphanum"`
	Currency     string  `db:"currency" json:"currency" validate:"required,iso4217"`
	Provider     string  `db:"provider" json:"provider" validate:"required"`
	Amount       float64 `db:"amount" json:"amount" validate:"required,gt=0"`
	PaymentDT    int64   `db:"payment_dt" json:"payment_dt" validate:"required"`
	Bank         string  `db:"bank" json:"bank" validate:"required"`
	DeliveryCost float64 `db:"delivery_cost" json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   float64 `db:"goods_total" json:"goods_total" validate:"required,gte=0"`
	CustomFee    float64 `db:"custom_fee" json:"custom_fee" validate:"gte=0"`
}

func (p *Payment) GetDataForSQLString() []any {
	return []any{
		p.Transaction, p.RequestID, p.Currency,
		p.Provider, p.Amount, p.PaymentDT,
		p.Bank, p.DeliveryCost, p.GoodsTotal,
		p.CustomFee,
	}
}

func (p *Payment) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	stringFields := []string{
		p.Transaction,
		p.RequestID,
		p.Currency,
		p.Provider,
		p.Bank,
	}
	for _, field := range stringFields {
		if err := binaryutils.WriteString(buf, field); err != nil {
			return nil, err
		}
	}

	floatFields := []float64{
		p.Amount,
		p.DeliveryCost,
		p.GoodsTotal,
		p.CustomFee,
	}
	for _, field := range floatFields {
		if err := binary.Write(buf, binary.LittleEndian, math.Float64bits(field)); err != nil {
			return nil, fmt.Errorf("failed to write float field: %w", err)
		}
	}

	if err := binary.Write(buf, binary.LittleEndian, p.PaymentDT); err != nil {
		return nil, fmt.Errorf("failed to write PaymentDT: %w", err)
	}

	return buf.Bytes(), nil
}

func (p *Payment) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	stringFields := []*string{
		&p.Transaction,
		&p.RequestID,
		&p.Currency,
		&p.Provider,
		&p.Bank,
	}
	for _, fieldPtr := range stringFields {
		str, err := binaryutils.ReadString(r)
		if err != nil {
			return err
		}
		*fieldPtr = str
	}

	floatFields := []*float64{
		&p.Amount,
		&p.DeliveryCost,
		&p.GoodsTotal,
		&p.CustomFee,
	}
	for _, fieldPtr := range floatFields {
		var bits uint64
		if err := binary.Read(r, binary.LittleEndian, &bits); err != nil {
			return fmt.Errorf("failed to read float field: %w", err)
		}
		*fieldPtr = math.Float64frombits(bits)
	}

	if err := binary.Read(r, binary.LittleEndian, &p.PaymentDT); err != nil {
		return fmt.Errorf("failed to read PaymentDT: %w", err)
	}

	return nil
}
