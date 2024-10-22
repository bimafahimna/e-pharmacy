create table "pharmacy_products" (
    "pharmacy_id" bigint not null,
    "product_id" bigint not null,
    "sold_amount" int not null,
    "price" decimal not null,
    "stock" int not null,
    "is_active" boolean not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    "deleted_at" timestamp null,
    constraint "pharmacy_products_pkey" primary key ("pharmacy_id", "product_id")
);

alter table "pharmacy_products"
add foreign key ("pharmacy_id") references "pharmacies" ("id");

alter table "pharmacy_products"
add foreign key ("product_id") references "products" ("id");