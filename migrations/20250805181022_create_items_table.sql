-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE items(
    chrt_id serial primary key,
    track_number text not null,
    price decimal(12, 2) not null,
    rid text not null,
    name varchar(255) not null,
    sale smallint not null,
    size varchar(10) not null,
    total_price decimal(12, 2) not null,
    nm_id bigint not null,
    brand varchar(255) not null,
    status smallint not null
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE items;
-- +goose StatementEnd
