alter table if exists "customer_addresses"
drop constraint "customer_addresses_user_id_fkey";

drop table if exists "customer_addresses";