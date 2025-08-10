package delivery

import (
	"bytes"
	binaryutils "first-task/pkg/utils/binaryUtils"
)

type Delivery struct {
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
