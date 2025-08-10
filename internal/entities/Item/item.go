package item

import (
	"bytes"
	"encoding/binary"
	binaryutils "first-task/pkg/utils/binaryUtils"
	"fmt"
	"math"
)

type Item struct {
	ChrtID      int64   `db:"chrt_id" json:"chrt_id"`
	TrackNumber string  `db:"track_number" json:"track_number"`
	Price       float64 `db:"price" json:"price"`
	RID         string  `db:"rid" json:"rid"`
	Name        string  `db:"name" json:"name"`
	Sale        uint8   `db:"sale" json:"sale"`
	Size        string  `db:"size" json:"size"`
	TotalPrice  float64 `db:"total_price" json:"total_price"`
	NMID        int64   `db:"nm_id" json:"nm_id"`
	Brand       string  `db:"brand" json:"brand"`
	Status      int32   `db:"status" json:"status"`
}

func (i *Item) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	strFields := []string{i.TrackNumber, i.RID, i.Name, i.Size, i.Brand}
	for _, v := range strFields {
		if err := binaryutils.WriteString(buf, v); err != nil {
			return nil, err
		}
	}

	if err := binary.Write(buf, binary.LittleEndian, i.ChrtID); err != nil {
		return nil, fmt.Errorf("failed to write ChrtID: %w", err)
	}

	if err := binary.Write(buf, binary.LittleEndian, math.Float64bits(i.Price)); err != nil {
		return nil, fmt.Errorf("failed to write Price: %w", err)
	}

	if err := binary.Write(buf, binary.LittleEndian, i.Sale); err != nil {
		return nil, fmt.Errorf("failed to write Sale: %w", err)
	}

	if err := binary.Write(buf, binary.LittleEndian, math.Float64bits(i.TotalPrice)); err != nil {
		return nil, fmt.Errorf("failed to write TotalPrice: %w", err)
	}

	if err := binary.Write(buf, binary.LittleEndian, i.NMID); err != nil {
		return nil, fmt.Errorf("failed to write NMID: %w", err)
	}

	if err := binary.Write(buf, binary.LittleEndian, i.Status); err != nil {
		return nil, fmt.Errorf("failed to write Status: %w", err)
	}

	return buf.Bytes(), nil
}

func (i *Item) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	strFields := []*string{&i.TrackNumber, &i.RID, &i.Name, &i.Size, &i.Brand}
	for _, v := range strFields {
		str, err := binaryutils.ReadString(r)
		if err != nil {
			return err
		}
		*v = str
	}

	if err := binary.Read(r, binary.LittleEndian, &i.ChrtID); err != nil {
		return fmt.Errorf("failed to read ChrtID: %w", err)
	}

	var priceBits uint64
	if err := binary.Read(r, binary.LittleEndian, &priceBits); err != nil {
		return fmt.Errorf("failed to read Price: %w", err)
	}
	i.Price = math.Float64frombits(priceBits)

	if err := binary.Read(r, binary.LittleEndian, &i.Sale); err != nil {
		return fmt.Errorf("failed to read Sale: %w", err)
	}

	var totalPriceBits uint64
	if err := binary.Read(r, binary.LittleEndian, &totalPriceBits); err != nil {
		return fmt.Errorf("failed to read TotalPrice: %w", err)
	}
	i.TotalPrice = math.Float64frombits(totalPriceBits)

	if err := binary.Read(r, binary.LittleEndian, &i.NMID); err != nil {
		return fmt.Errorf("failed to read NMID: %w", err)
	}

	if err := binary.Read(r, binary.LittleEndian, &i.Status); err != nil {
		return fmt.Errorf("failed to read Status: %w", err)
	}

	return nil
}
