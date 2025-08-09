package storage

import (
	order "first-task/internal/entities/Order"
	"fmt"

	"go.uber.org/zap"
)

func (s *Storage) LoadInitialData(size int) {
	zap.L().Info("start initialization cache")

	ords, err := s.dataBaseStorage.GetInitialData(size)
	if err != nil {
		zap.L().Error("cant load initial cache")
	}

	s.localStorage.LoadInitialCache(ords)
}

func (s *Storage) AddOrder(ord *order.Order) error {
	const op = "internal.storage.AddOrder"

	err := s.dataBaseStorage.Add(ord)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) FindOrder(orderUID string) (*order.Order, error) {
	const op = "internal.storage.FindOrder"

	var result *order.Order
	var err error

	result = s.localStorage.Find(orderUID)
	if result != nil {
		return result, nil
	}

	result, err = s.dataBaseStorage.Find(orderUID)
	if err != nil {
		return &order.Order{}, fmt.Errorf("%s: %w", op, err)
	}
	s.localStorage.Add(result)

	return result, nil
}
