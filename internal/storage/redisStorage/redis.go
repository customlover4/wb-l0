package redisStorage

import (
	"first-task/internal/config"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	Host     = "localhost"
	Port     = "6379"
	Password = ""
	DB       = 0
)

type RedisStorage struct {
	rdb *redis.Client
}

func NewRedisStorage(cfg config.RedisConfig) RedisStorage {
	return RedisStorage{
		redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DBName,
		}),
	}
}
