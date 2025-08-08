package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var GetOrderJSONFromDataBase = fmt.Sprintf(`
	select json_build_object (
		'order_uid', o.order_uid, 
		'track_number', o.track_number, 
		'entry', o.entry, 
		'locale', o.locale, 
		'internal_signature', o.internal_signature, 
		'customer_id', o.customer_id, 
		'delivery_service', o.delivery_service, 
		'shardkey', o.shardkey, 
		'sm_id', o.sm_id, 
		'date_created', o.date_created,
		'oof_shard', o.oof_shard,
		'delivery', json_build_object(
			'id', di.id,
			'name', di.name,
			'phone', di.phone, 
			'zip', di.zip,
			'city', di.city,
			'address', di.address,
			'region', di.region, 
			'email', di.email
		),
		'payment', json_build_object(
			'id', p.id,
			'transaction', p.transaction,
			'request_id', p.request_id,
			'currency', p.currency,
			'provider', p.provider,
			'amount', p.amount,
			'payment_dt', p.payment_dt,
			'bank', p.bank,
			'delivery_cost', p.delivery_cost,
			'goods_total', p.goods_total,
			'custom_fee', p.custom_fee
		),
		'items', (
			select json_agg(json_build_object(
				'chrt_id', i.chrt_id,
				'track_number', i.track_number,
				'price', i.price,
				'rid', i.rid,
				'name', i.name,
				'sale', i.sale,
				'size', i.size,
				'total_price', i.total_price,
				'nm_id', i.nm_id,
				'brand', i.brand,
				'status', i.status
			))
			from orders_items as oi join items as i on oi.item_id=i.chrt_id
			where oi.order_id=o.id
		)
	)
	from %s as o join %s as di on o.delivery_id=di.id 
	join %s as p on o.payment_id=p.id
	where order_uid=$1;
`, OrdersTable, DeliveryInfoTable, PaymentInfoTable)

const InitialRequestLength = 100

func GetLastOrdersJSONFromDataBase(size int) string {
	return fmt.Sprintf(`
select json_build_object (
		'order_uid', o.order_uid, 
		'track_number', o.track_number, 
		'entry', o.entry, 
		'locale', o.locale, 
		'internal_signature', o.internal_signature, 
		'customer_id', o.customer_id, 
		'delivery_service', o.delivery_service, 
		'shardkey', o.shardkey, 
		'sm_id', o.sm_id, 
		'date_created', o.date_created,
		'oof_shard', o.oof_shard,
		'delivery', json_build_object(
			'id', di.id,
			'name', di.name,
			'phone', di.phone, 
			'zip', di.zip,
			'city', di.city,
			'address', di.address,
			'region', di.region, 
			'email', di.email
		),
		'payment', json_build_object(
			'id', p.id,
			'transaction', p.transaction,
			'request_id', p.request_id,
			'currency', p.currency,
			'provider', p.provider,
			'amount', p.amount,
			'payment_dt', p.payment_dt,
			'bank', p.bank,
			'delivery_cost', p.delivery_cost,
			'goods_total', p.goods_total,
			'custom_fee', p.custom_fee
		),
		'items', (
			select json_agg(json_build_object(
				'chrt_id', i.chrt_id,
				'track_number', i.track_number,
				'price', i.price,
				'rid', i.rid,
				'name', i.name,
				'sale', i.sale,
				'size', i.size,
				'total_price', i.total_price,
				'nm_id', i.nm_id,
				'brand', i.brand,
				'status', i.status
			))
			from orders_items as oi join items as i on oi.item_id=i.chrt_id
			where oi.order_id=o.id
		)
	)
	from %s as o join %s as di on o.delivery_id=di.id 
	join %s as p on o.payment_id=p.id order by o.id desc limit %d;
`, OrdersTable, DeliveryInfoTable, PaymentInfoTable, size)
}

func GetInsertPaymentSQLString() string {
	return fmt.Sprintf(`
	insert into %s (
	transaction, request_id, currency, provider, amount, payment_dt, bank, 
	delivery_cost, goods_total, custom_fee
	) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id;
	`, PaymentInfoTable)
}

func GetInsertOrderSQLString() string {
	return fmt.Sprintf(`
	insert into %s (order_uid, track_number, entry, delivery_id, payment_id,
	locale, internal_signature, customer_id, delivery_service, shardkey,
	sm_id, date_created, oof_shard)
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) 
	returning id;`, OrdersTable)
}

func GetInsertDeliverySQLString() string {
	return fmt.Sprintf(`
	insert into %s (name, phone, zip, city, address, region, email) 
	values (
	$1, $2, $3, $4, $5, $6, $7
	) returning id;
	`, DeliveryInfoTable)
}

func GetInsertOrdersItemsSQLString(itmsLen int) string {
	res := fmt.Sprintf(
		"insert into %s (order_id, item_id) values ", OrdersItemsTable,
	)

	for j := 1; j <= itmsLen*2; j += 2 {
		if j > 1 {
			res += ", "
		}
		res += fmt.Sprintf("($%d, $%d)", j, j+1)
	}
	res += ";"

	return res
}

func HandleTxErr(tx *sqlx.Tx, InternalErr error) error {
	if err := tx.Rollback(); err != nil {
		return err
	}
	return InternalErr
}
