-- +goose Up
-- +goose StatementBegin
create table orders
(
    id      bigserial primary key,
    user_id bigint not null,
    status  int    not null
);
-- +goose StatementEnd

-- +goose StatementBegin
create table orders_items
(
    id       bigserial primary key,
    order_id bigint not null,
    sku      bigint not null,
    count    int    not null
);
-- +goose StatementEnd

-- +goose StatementBegin
create table stocks
(
    id          bigserial primary key,
    sku         bigint not null unique,
    total_count int    not null,
    reserved    int    not null
);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
drop table orders;
drop table orders_items;
drop table stocks;
-- +goose StatementEnd
