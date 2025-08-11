package integrational

import (
	"context"
	"database/sql"
	"first-task/internal/config"
	"first-task/internal/storage/postgres"
	"first-task/internal/storage/redisStorage"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DBHost     = "localhost"
	DBPort     = "5432"
	DBUser     = "test"
	DBPassword = "test"
	DBName     = "testdb"
	KafkaTopic = "orders"
)

func SetupTestDB(t *testing.T) (testcontainers.Container, *postgres.Postgres) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     DBUser,
			"POSTGRES_PASSWORD": DBPassword,
			"POSTGRES_DB":       DBName,
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, _ := pgContainer.Host(ctx)
	port, _ := pgContainer.MappedPort(ctx, "5432")
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port.Port(), DBUser, DBPassword, DBName,
	)

	time.Sleep(time.Second * 3)

	oldDB, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	defer oldDB.Close()

	err = oldDB.Ping()
	require.NoError(t, err)

	applyMigrations(t, oldDB)

	pg := postgres.NewPostgres(config.PostgresConfig{
		Host:     host,
		Port:     port.Port(),
		User:     DBUser,
		Password: DBPassword,
		DBName:   DBName,
		SSLMode:  false,
	})

	return pgContainer, pg
}

func applyMigrations(t *testing.T, db *sql.DB) {
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "../../../migrations")

	goose.SetBaseFS(nil)

	if err := goose.Up(db, migrationsDir); err != nil {
		t.Fatalf("Failed to apply migrations: %v", err)
	}

	_, err := db.Exec(MockItemRow)
	require.NoError(t, err)
}

var MockItemRow = `
insert into items (
	chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, 
	brand, status
)
values(
	9934930, 'WBILMTESTTRACK', 453, 'ab4219087a764ae0btest', 'Mascaras',
	30, '0', 317, 2389212, 'Vivienne Sabo', 202
);
`

func SetupTestRedis(t *testing.T) (testcontainers.Container, *redisStorage.RedisStorage) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7.2-alpine",
		ExposedPorts: []string{"6379/tcp"},
		Env: map[string]string{
			"MAXMEMORY":        "100MB",
			"MAXMEMORY_POLICY": "volatile-ttl",
		},
		WaitingFor: wait.ForLog("Ready to accept connections"),
	}
	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	time.Sleep(time.Second * 3)

	host, err := redisContainer.Host(ctx)
	require.NoError(t, err)
	port, err := redisContainer.MappedPort(ctx, "6379")
	require.NoError(t, err)

	str := redisStorage.NewRedisStorage(config.RedisConfig{
		Host: host,
		Port: port.Port(),
	})

	return redisContainer, str
}

func SetupTestKafka(t *testing.T) testcontainers.Container {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "bitnami/kafka:4.0.0",
		ExposedPorts: []string{"9092/tcp"},
		Env: map[string]string{
			"KAFKA_CFG_NODE_ID":                        "1",
			"KAFKA_CFG_PROCESS_ROLES":                  "controller,broker",
			"KAFKA_CFG_LISTENERS":                      "PLAINTEXT://:9092,CONTROLLER://:9093",
			"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP": "CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
			"KAFKA_CFG_ADVERTISED_LISTENERS":           "PLAINTEXT://localhost:9092",
			"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS":       "1@localhost:9093",
			"KAFKA_CFG_CONTROLLER_LISTENER_NAMES":      "CONTROLLER",
			"KAFKA_CFG_INTER_BROKER_LISTENER_NAME":     "PLAINTEXT",
		},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"9092/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0", // Слушаем на всех интерфейсах
						HostPort: "9092",    // Фиксируем порт хоста
					},
				},
			}
		},
		WaitingFor: wait.ForListeningPort("9092/tcp"),
	}
	kafkaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	return kafkaContainer
}
