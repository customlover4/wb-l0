-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE orders (
    id serial primary key,
    order_uid text not null unique,
    track_number text not null,
    entry varchar(255) not null,
    delivery_id bigint not null,
    FOREIGN KEY (delivery_id) REFERENCES delivery_info(id) ON DELETE CASCADE,
    payment_id bigint not null,
    FOREIGN KEY (payment_id) REFERENCES payment_info(id) ON DELETE CASCADE,
    -- many to many on items
    locale varchar(10) not null,
    internal_signature text not null,
    customer_id text not null,
    delivery_service varchar(255) not null,
    shardkey text not null,
    sm_id smallint not null,
    date_created varchar(255) not null,
    oof_shard text not null
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE orders;
-- +goose StatementEnd
