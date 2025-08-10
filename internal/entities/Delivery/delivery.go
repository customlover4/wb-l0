package delivery

import (
	"bytes"
	"encoding/binary"
	binaryutils "first-task/pkg/utils/binaryUtils"
	"fmt"
)

type Delivery struct {
	Id      int64  `db:"id" json:"-"`
	Name    string `db:"name" json:"name"`
	Phone   string `db:"phone" json:"phone"`
	Zip     string `db:"zip" json:"zip"`
	City    string `db:"city" json:"city"`
	Address string `db:"address" json:"address"`
	Region  string `db:"region" json:"region"`
	Email   string `db:"email" json:"email"`
}

func (d *Delivery) GetDataForSQLString() []any {
	return []any{
		d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email,
	}
}

func (d *Delivery) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, uint64(d.Id)); err != nil {
		return nil, fmt.Errorf("failed to write Id: %w", err)
	}

	stringFields := []string{d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email}
	for _, field := range stringFields {
		if err := binaryutils.WriteString(buf, field); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (d *Delivery) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)

	var id uint64
	if err := binary.Read(r, binary.LittleEndian, &id); err != nil {
		return fmt.Errorf("failed to read Id: %w", err)
	}
	d.Id = int64(id)

	stringFields := []*string{&d.Name, &d.Phone, &d.Zip, &d.City, &d.Address, &d.Region, &d.Email}
	for _, fieldPtr := range stringFields {
		r, err := binaryutils.ReadString(r)
		if err != nil {
			return err
		}
		*fieldPtr = r
	}

	return nil
}
