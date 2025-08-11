package integrational

import (
	"context"
	"database/sql"
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
	DBUser     = "test"
	DBPassword = "test"
	DBName     = "testdb"

	KafkaTopic = "test_topic"

	DBMapped    = "5432"
	RedisMapped = "6379"
	KafkaMapped = "9092"
)

func SetupTestDB(t *testing.T) testcontainers.Container {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", DBMapped)},
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
	port, _ := pgContainer.MappedPort(ctx, DBMapped)
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

	return pgContainer
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

func SetupTestRedis(t *testing.T) testcontainers.Container {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "redis:7.2-alpine",
		ExposedPorts: []string{fmt.Sprintf("%s/tcp", RedisMapped)},
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

	return redisContainer
}

func SetupTestKafka(t *testing.T) testcontainers.Container {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:           "bitnami/kafka:4.0.0",
		HostAccessPorts: []int{9092},
		ExposedPorts:    []string{"9092/tcp"},
		Env: map[string]string{
			"KAFKA_CFG_NODE_ID":       "1",
			"KAFKA_CFG_PROCESS_ROLES": "controller,broker",
			"KAFKA_CFG_LISTENERS": fmt.Sprintf(
				"PLAINTEXT://:%s,CONTROLLER://:9093", KafkaMapped,
			),
			"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP": "CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT",
			"KAFKA_CFG_ADVERTISED_LISTENERS": fmt.Sprintf(
				"PLAINTEXT://localhost:%s", KafkaMapped,
			),
			"KAFKA_CFG_CONTROLLER_QUORUM_VOTERS":   "1@localhost:9093",
			"KAFKA_CFG_CONTROLLER_LISTENER_NAMES":  "CONTROLLER",
			"KAFKA_CFG_INTER_BROKER_LISTENER_NAME": "PLAINTEXT",
		},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"9092/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",   // Слушаем на всех интерфейсах
						HostPort: KafkaMapped, // Фиксируем порт хоста
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
