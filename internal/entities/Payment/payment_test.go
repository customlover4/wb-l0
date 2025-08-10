package payment

import "testing"

func CmpPayment(a, b *Payment) bool {
	return a.Transaction == b.Transaction &&
		a.RequestID == b.RequestID &&
		a.Currency == b.Currency &&
		a.Provider == b.Provider &&
		a.Amount == b.Amount &&
		a.PaymentDT == b.PaymentDT &&
		a.Bank == b.Bank &&
		a.DeliveryCost == b.DeliveryCost &&
		a.GoodsTotal == b.GoodsTotal &&
		a.CustomFee == b.CustomFee
}

func TestPaymentMarshaling(t *testing.T) {
	p := &Payment{
		Transaction:  "test_transaction",
		RequestID:    "test_request",
		Currency:     "USD",
		Provider:     "test_provider",
		Amount:       1000,
		PaymentDT:    1234567890,
		Bank:         "test_bank",
		DeliveryCost: 500,
		GoodsTotal:   1500,
		CustomFee:    0,
	}

	tmp, err := p.MarshalBinary()
	if err != nil {
		t.Error("failed on marshaling Payment: " + err.Error())
		return
	}

	newP := &Payment{}
	err = newP.UnmarshalBinary(tmp)
	if err != nil {
		t.Error("failed on unmarshaling Payment: " + err.Error())
		return
	}

	if !CmpPayment(p, newP) {
		t.Errorf("Payment unmarshaling failed, fields don't match")
		t.Errorf("Original: %+v", p)
		t.Errorf("Unmarshaled: %+v", newP)
	}
}
