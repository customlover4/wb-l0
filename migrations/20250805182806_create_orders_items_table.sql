-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE orders_items(
    id serial primary key,
    order_id bigint not null,
    item_id bigint not null,
    foreign key (order_id) references orders(id) on delete cascade,
    foreign key (item_id) references items(chrt_id) on delete cascade
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
drop table orders_items;
-- +goose StatementEnd
