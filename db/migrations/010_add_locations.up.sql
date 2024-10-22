create table "provinces" (
    "id" bigserial not null,
    "name" varchar not null,
    constraint "provinces_pkey" primary key ("id")
);

create table "cities" (
    "id" bigserial not null,
    "province_id" bigint not null,
    "name" varchar not null,
    "type" varchar not null,
    "unofficial_id" bigint null,
    "province_unofficial_id" bigint not null,
    constraint "cities_pkey" primary key ("id")
);

create table "districts" (
    "id" bigserial not null,
    "city_id" bigint not null,
    "name" varchar not null,
    constraint "districts_pkey" primary key ("id")
);

create table "sub_districts" (
    "id" bigserial not null,
    "district_id" bigint not null,
    "name" varchar not null,
    constraint "sub_districts_pkey" primary key ("id")
);

alter table "cities"
add foreign key ("province_unofficial_id") references "provinces" ("id");

alter table "districts"
add foreign key ("city_id") references "cities" ("id");

alter table "sub_districts"
add foreign key ("district_id") references "districts" ("id");
