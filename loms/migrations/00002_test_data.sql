-- +goose Up
-- +goose StatementBegin
insert into stocks (sku, total_count, reserved) values (139275865, 65534, 0);
insert into stocks (sku, total_count, reserved) values (1076963, 65534, 10);
insert into stocks (sku, total_count, reserved) values (1148162, 200, 20);
insert into stocks (sku, total_count, reserved) values (1625903, 250, 30);
insert into stocks (sku, total_count, reserved) values (2618151, 300, 40);
insert into stocks (sku, total_count, reserved) values (2956315, 350, 50);

insert into orders (user_id, status) values (1, 0); -- NEW
insert into orders_items (order_id, sku, count) values (1, 1076963, 5);

insert into orders (user_id, status) values (2, 1); -- AWAITING_PAYMENT
insert into orders_items (order_id, sku, count) values (2, 1076963, 5);

insert into orders (user_id, status) values (3, 2); -- FAILED
insert into orders_items (order_id, sku, count) values (3, 1148162, 10);

insert into orders (user_id, status) values (4, 3); -- PAYED
insert into orders_items (order_id, sku, count) values (4, 1148162, 10);

insert into orders (user_id, status) values (5, 4); -- CANCELLED
insert into orders_items (order_id, sku, count) values (5, 1148162, 10);

-- +goose StatementEnd