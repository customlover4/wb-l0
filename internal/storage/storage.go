package storage

import (
	"errors"
	order "first-task/internal/entities/Order"
)

type Storage struct {
	localStorage    Cacher
	dataBaseStorage DataBaser
}

var ErrNotFound = errors.New("not found")

func NewStorage(ls Cacher, dbs DataBaser) *Storage {
	if ls == nil || dbs == nil {
		panic("can't create storage without one or two storagers")
	}
	return &Storage{
		localStorage:    ls,
		dataBaseStorage: dbs,
	}
}

type DataBaser interface {
	Add(ord *order.Order) error
	Find(orderUID string) (*order.Order, error)
	GetInitialData(size int) ([]*order.Order, error)
	Shutdown()
}

type Cacher interface {
	Add(ord *order.Order)
	Find(orderUID string) *order.Order
	Delete(orderUID string)
	LoadInitialCache(ords []*order.Order)
	Shutdown()
}
