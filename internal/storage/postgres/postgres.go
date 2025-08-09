package postgres

import (
	"errors"
	"first-task/internal/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	DeliveryInfoTable = "delivery_info"
	PaymentInfoTable  = "payment_info"
	ItemsTable        = "items"
	OrdersTable       = "orders"
	OrdersItemsTable  = "orders_items"
)

var (
	ErrNotFound = errors.New("can't find this order or wrong order_uid")
)

type Postgres struct {
	conn *sqlx.DB
}

func NewPostgres(cp config.PostgresConfig) *Postgres {
	connString := fmt.Sprintf(
		`host=%s port=%s user=%s password=%s dbname=%s sslmode=disable`,
		cp.Host, cp.Port, cp.User, cp.Password, cp.DBName,
	)

	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return &Postgres{db}
}

func (p Postgres) Shutdown() {
	if err := p.conn.Close(); err != nil {
		zap.L().Error(err.Error())
	}
}
