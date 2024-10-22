alter table if exists "pharmacy_products" drop constraint "pharmacy_products_product_id_fkey";

alter table if exists "pharmacy_products" drop constraint "pharmacy_products_pharmacy_id_fkey";

drop table if exists "pharmacy_products";