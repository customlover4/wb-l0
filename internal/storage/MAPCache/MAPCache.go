package mapcache

import (
	order "first-task/internal/entities/Order"
	"time"

	"go.uber.org/zap"
)

type MAPData struct {
	Expired int64
	Data    *order.Order
}

type MAPStorage struct {
	Cache map[string]MAPData
}

func NewMapStorage() *MAPStorage {
	r := &MAPStorage{
		make(map[string]MAPData),
	}
	go func() {
		ticker := time.NewTicker(time.Minute * 1)

		for dt := range ticker.C {
			zap.L().Info("checking cache for expired values " + dt.String())
			go r.clean()
		}
	}()
	return r
}

func (ms MAPStorage) Shutdown() {

}
