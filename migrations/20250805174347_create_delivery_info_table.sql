-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE delivery_info(
    id serial primary key,
    name varchar(255) not null,
    phone varchar(100) not null,
    zip varchar(100) not null,
    city varchar(255) not null,
    address varchar(255) not null,
    region varchar(255) not null,
    email varchar(255) not null
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE delivery_info;