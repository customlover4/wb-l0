package mapcache

import (
	order "first-task/internal/entities/Order"
	"time"
)

func (s *MAPStorage) clean() {
	stamp := time.Now().Unix()
	for _, v := range s.Cache {
		if (v.Expired - stamp) <= 0 {
			delete(s.Cache, v.Data.OrderUID)
		}
	}
}

func (s *MAPStorage) LoadInitialCache(ords []*order.Order) {
	for _, v := range ords {
		if v == nil {
			continue
		}
		s.Add(v)
	}
}

func (s *MAPStorage) Add(ord *order.Order) {
	if _, ok := s.Cache[ord.OrderUID]; !ok {
		s.Cache[ord.OrderUID] = MAPData{
			Expired: time.Now().Add(time.Second * 10).Unix(),
			Data:    ord,
		}
	}
}

func (s *MAPStorage) Find(orderUID string) *order.Order {
	if v, ok := s.Cache[orderUID]; ok {
		tmp := s.Cache[orderUID]
		tmp.Expired = time.Now().Add(time.Second * 10).Unix()
		s.Cache[orderUID] = tmp
		return v.Data
	}

	return nil
}

func (s *MAPStorage) Delete(orderUID string) {
	delete(s.Cache, orderUID)
}
