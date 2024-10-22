alter table if exists "products" drop constraint "products_name_generic_name_manufacturer_id_key";

alter table if exists "products" drop constraint "products_product_form_id_fkey";

alter table if exists "products" drop constraint "products_product_classification_id_fkey";

alter table if exists "products" drop constraint "products_manufacturer_id_fkey";

drop table if exists "product_forms";

drop table if exists "product_classifications";

drop table if exists "manufacturers";

drop table if exists "categories";

drop table if exists "products";