alter table if exists "order_items" drop constraint "order_items_product_id_fkey";

alter table if exists "order_items" drop constraint "order_items_order_id_fkey";

alter table if exists "orders" drop constraint "orders_logistic_id_fkey";

alter table if exists "orders" drop constraint "orders_payment_id_fkey";

alter table if exists "orders" drop constraint "orders_pharmacy_id_fkey";

alter table if exists "orders" drop constraint "orders_user_id_fkey";

drop table if exists "payments";

drop table if exists "order_items";

drop table if exists "orders";
