alter table if exists "cities"
drop constraint "cities_province_unofficial_id_fkey";

alter table if exists "districts"
drop constraint "districts_city_id_fkey";

alter table if exists "sub_districts"
drop constraint "sub_districts_district_id_fkey";

drop table if exists "sub_districts";

drop table if exists "districts";

drop table if exists "provinces";

drop table if exists "cities";
