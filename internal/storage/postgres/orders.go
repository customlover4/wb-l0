package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	order "first-task/internal/entities/Order"
	"first-task/internal/storage"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func (p *Postgres) Add(ord *order.Order) error {
	const op = "internal.storage.postgres.AddOrder"

	var lastInsertDeliverID int64
	var lastInsertPaymentID int64
	var lastInsertIDOrder int64

	transaction, err := p.conn.Beginx()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = sqlx.Get(
		transaction, &lastInsertDeliverID, GetInsertDeliverySQLString(),
		ord.Delivery.GetDataForSQLString()...,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleTxErr(transaction, err))
	}

	err = sqlx.Get(
		transaction, &lastInsertPaymentID, GetInsertPaymentSQLString(),
		ord.Payment.GetDataForSQLString()...,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleTxErr(transaction, err))
	}

	err = sqlx.Get(
		transaction, &lastInsertIDOrder, GetInsertOrderSQLString(),
		ord.GetDataForSQLString(lastInsertDeliverID, lastInsertPaymentID)...,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleTxErr(transaction, err))
	}

	_, err = transaction.Exec(
		GetInsertOrdersItemsSQLString(len(ord.Items)),
		ord.GetDataForSQLStringOrdersItems(lastInsertIDOrder)...,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleTxErr(transaction, err))
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, HandleTxErr(transaction, err))
	}

	return nil
}

func (p *Postgres) Find(orderUID string) (*order.Order, error) {
	const op = "internal.storage.postgres.FindOrder"

	var tmp []byte
	err := p.conn.Get(&tmp, GetOrderJSONFromDataBase, orderUID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, storage.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var result *order.Order
	err = json.Unmarshal(tmp, &result)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (p *Postgres) GetInitialData(size int) ([]*order.Order, error) {
	const op = "internal.storage.postgres.GetInitialData"

	tmpData := make([]string, 0, size)
	result := make([]*order.Order, size)

	err := p.conn.Select(&tmpData, GetLastOrdersJSONFromDataBase(size))
	if err != nil {
		return []*order.Order{}, fmt.Errorf("%s: %w", op, err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(tmpData))
	for i, v := range tmpData {
		go func(wg *sync.WaitGroup, i int, str string) {
			defer wg.Done()
			newOrd := new(order.Order)
			if err := json.Unmarshal([]byte(str), newOrd); err == nil {
				result[i] = newOrd
			}
		}(wg, i, v)
	}
	wg.Wait()

	filteredResult := make([]*order.Order, 0, size)
	for _, v := range result {
		if v == nil {
			continue
		}
		filteredResult = append(filteredResult, v)
	}

	return filteredResult, nil
}
