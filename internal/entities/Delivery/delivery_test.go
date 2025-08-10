package delivery

import "testing"

func CmpDelivery(a, b *Delivery) bool {
	return a.Name == b.Name &&
		a.Phone == b.Phone &&
		a.Zip == b.Zip &&
		a.City == b.City &&
		a.Address == b.Address &&
		a.Region == b.Region &&
		a.Email == b.Email
}

func TestMarshaling(t *testing.T) {
	d := &Delivery{
		Name:    "test",
		Phone:   "test",
		Zip:     "test",
		City:    "test",
		Address: "test",
		Region:  "test",
		Email:   "test",
	}

	tmp, err := d.MarshalBinary()
	if err != nil {
		t.Error("failed on marshaling" + err.Error())
		return
	}

	newD := &Delivery{}
	err = newD.UnmarshalBinary(tmp)
	if err != nil {
		t.Error("failed on marshaling" + err.Error())
		return
	}

	if !CmpDelivery(d, newD) {
		t.Errorf("Item unmarshaling failed, fields don't match")
		t.Errorf("Original: %+v", d)
		t.Errorf("Unmarshaled: %+v", newD)
	}
}
