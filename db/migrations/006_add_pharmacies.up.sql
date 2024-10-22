create table "pharmacies" (
    "id" bigserial not null,
    "pharmacist_id" bigint null,
    "partner_id" bigint not null,
    "name" varchar not null,
    "address" varchar not null,
    "city_id" bigint not null,
    "latitude" decimal not null,
    "longitude" decimal not null,
    "is_active" boolean not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    "deleted_at" timestamp null,
    constraint "pharmacies_pkey" primary key ("id"),
    unique ("latitude", "longitude")
);

create table "partners" (
    "id" bigserial not null,
    "name" varchar not null unique,
    "logo_url" varchar not null,
    "year_founded" int not null,
    "active_days" varchar not null,
    "operational_start" varchar not null,
    "operational_stop" varchar not null,
    "is_active" boolean not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    "deleted_at" timestamp null,
    constraint "partners_pkey" primary key ("id")
);

create table "logistics" (
    "id" bigserial not null,
    "name" varchar not null,
    "logo_url" varchar not null,
    "service" varchar not null,
    "price_per_kilometers" decimal null,
    "eda" int not null,
    constraint "logistics_pkey" primary key ("id"),
    unique ("name", "service")
);

create table "pharmacy_logistics" (
    "pharmacy_id" bigint not null,
    "logistic_id" bigint not null,
    constraint "pharmacy_logistics_pkey" primary key ("pharmacy_id", "logistic_id")
);

alter table "pharmacies"
add foreign key ("pharmacist_id") references "users" ("id");

alter table "pharmacies"
add foreign key ("partner_id") references "partners" ("id");

alter table "pharmacy_logistics"
add foreign key ("pharmacy_id") references "pharmacies" ("id");

alter table "pharmacy_logistics"
add foreign key ("logistic_id") references "logistics" ("id");