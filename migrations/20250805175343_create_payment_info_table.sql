-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE payment_info(
    id serial primary key,
    transaction varchar(255) not null,
    request_id text not null unique,
    currency varchar(20) not null,
    provider varchar(255) not null,
    amount decimal(12, 2) not null,
    payment_dt bigint not null,
    bank varchar(255) not null,
    delivery_cost decimal(12, 2) not null,
    goods_total decimal(12, 2) not null,
    custom_fee decimal(12, 2) not null
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table payment_info;
-- +goose StatementEnd
