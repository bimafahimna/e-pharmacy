alter table if exists "pharmacist_details" drop constraint "pharmacist_details_user_id_fkey";

alter table if exists "customer_details" drop constraint "customer_details_user_id_fkey";

drop table if exists "pharmacist_details";

drop table if exists "customer_details";