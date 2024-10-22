alter table if exists "carts" drop constraint "carts_product_id_fkey";

alter table if exists "carts" drop constraint "carts_pharmacy_id_fkey";

alter table if exists "carts" drop constraint "carts_user_id_fkey";

drop table if exists "carts";