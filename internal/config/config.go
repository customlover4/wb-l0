package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	WebConfig         `yaml:"web_config" env-required:"true"`
	PostgresConfig    `yaml:"postgres_config"`
	RedisConfig       `yaml:"redis"`
	KafkaOrdersConfig `yaml:"kafka"`

	InitialDataSize int `yaml:"initial_data_size" env-default:"100"`
}

type WebConfig struct {
	Host         string        `yaml:"host" env-required:"true"`
	Port         string        `yaml:"port" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"10s"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	SSLMode  bool   `yaml:"sslmode" env-default:"false"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	Password string `yaml:"password"`
	DBName   int    `yaml:"db_name"`
}

type KafkaOrdersConfig struct {
	Brokers  []string `yaml:"brokers" env-required:"true"`
	Topic    string   `yaml:"topic" env-required:"true"`
	MinBytes int      `yaml:"min_bytes" env-default:"1"`
	MaxBytes int      `yaml:"max_bytes" env-default:"10e6"`
	GroupID  string   `yaml:"group_id" env-required:"true"`
}

// if can't find config file throw panic
func MustLoad(filePath string) *Config {
	f, err := os.Open(filePath)
	defer func() {
		err := f.Close()
		if err != nil {
			zap.L().Panic(err.Error())
		}
	}()
	if err != nil {
		zap.L().Panic("can't find config file")
	}

	cfg := new(Config)
	if err := cleanenv.ReadConfig(filePath, cfg); err != nil {
		zap.L().Panic("error on reading config file")
	}

	return cfg
}
