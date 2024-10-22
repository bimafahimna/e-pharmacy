alter table if exists "pharmacy_logistics" drop constraint "pharmacy_logistics_logistic_id_fkey";

alter table if exists "pharmacy_logistics" drop constraint "pharmacy_logistics_pharmacy_id_fkey";

alter table if exists "pharmacies" drop constraint "pharmacies_partner_id_fkey";

alter table if exists "pharmacies" drop constraint "pharmacies_pharmacist_id_fkey";

drop table if exists "pharmacy_logistics";

drop table if exists "logistics";

drop table if exists "partners";

drop table if exists "pharmacies";