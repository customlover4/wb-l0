package redisStorage

import (
	"context"
	order "first-task/internal/entities/Order"
	"time"

	"go.uber.org/zap"
)

func (rs *RedisStorage) LoadInitialCache(ords []*order.Order) {
	for _, v := range ords {
		rs.Add(v)
	}
}

func (rs *RedisStorage) Add(ord *order.Order) {
	res := rs.rdb.Set(context.Background(), ord.OrderUID, ord, time.Hour*1)
	if res.Err() != nil {
		zap.L().Error("on adding value to redis storage")
	}
}

func (rs *RedisStorage) Find(orderUID string) *order.Order {
	res := rs.rdb.Get(context.Background(), orderUID)
	if res.Err() != nil {
		return nil
	}

	var resultData order.Order
	err := res.Scan(&resultData)
	if err != nil {
		zap.L().Error("on scanning value from redis storage")
		return nil
	}

	return &resultData
}

func (rs *RedisStorage) Delete(orderUID string) {
	r := rs.rdb.Del(context.Background(), orderUID)
	if r.Err() != nil {
		zap.L().Error("on deleting value from redis storage")
	}
}
