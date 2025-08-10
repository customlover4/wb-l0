package item

import "testing"

func CmpItem(a, b *Item) bool {
	return a.ChrtID == b.ChrtID &&
		a.TrackNumber == b.TrackNumber &&
		a.Price == b.Price &&
		a.RID == b.RID &&
		a.Name == b.Name &&
		a.Sale == b.Sale &&
		a.Size == b.Size &&
		a.TotalPrice == b.TotalPrice &&
		a.NMID == b.NMID &&
		a.Brand == b.Brand &&
		a.Status == b.Status
}

func TestItemMarshaling(t *testing.T) {
	i := &Item{
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
	}

	tmp, err := i.MarshalBinary()
	if err != nil {
		t.Error("failed on marshaling Item: " + err.Error())
		return
	}

	newI := &Item{}
	err = newI.UnmarshalBinary(tmp)
	if err != nil {
		t.Error("failed on unmarshaling Item: " + err.Error())
		return
	}

	if !CmpItem(i, newI) {
		t.Errorf("Item unmarshaling failed, fields don't match")
		t.Errorf("Original: %+v", i)
		t.Errorf("Unmarshaled: %+v", newI)
	}
}
