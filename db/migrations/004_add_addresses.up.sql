create table "customer_addresses" (
    "id" bigserial not null,
    "user_id" bigint not null,
    "name" varchar not null,
    "receiver_name" varchar not null,
    "receiver_phone_number" varchar not null,
    "latitude" decimal not null,
    "longitude" decimal not null,
    "province" varchar not null,
    "city_id" bigint not null,
    "city" varchar not null,
    "district" varchar not null,
    "sub_district" varchar not null,
    "address_details" varchar not null,
    "is_active" boolean not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    constraint "customer_addresses_pkey" primary key ("id"),
    unique ("user_id", "name")
);

alter table "customer_addresses"
add foreign key ("user_id") references "users" ("id");