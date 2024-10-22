create table "carts" (
    "user_id" bigint not null,
    "pharmacy_id" bigint not null,
    "product_id" bigint not null,
    "quantity" bigint not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    constraint "carts_pkey" primary key ("user_id", "pharmacy_id", "product_id")
);

alter table "carts"
add foreign key ("user_id") references "users" ("id");

alter table "carts"
add foreign key ("pharmacy_id") references "pharmacies" ("id");

alter table "carts"
add foreign key ("product_id") references "products" ("id");