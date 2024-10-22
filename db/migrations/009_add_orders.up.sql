create table "orders" (
    "id" bigserial not null,
    "user_id" bigint not null,
    "pharmacy_id" bigint not null,
    "payment_id" bigint not null,
    "pharmacy_name" varchar not null,
    "status" varchar not null,
    "address" text not null,
    "contact_name" varchar not null,
    "contact_phone" varchar not null,
    "logistic_id" bigint not null,
    "logistic_cost" decimal not null,
    "amount" decimal not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    constraint "orders_pkey" primary key ("id")
);

create table "order_items" (
    "id" bigserial not null,
    "order_id" bigint not null,
    "pharmacy_id" bigint not null,
    "product_id" bigint not null,
    "image_url" varchar not null,
    "name" varchar not null,
    "quantity" int not null,
    "price" decimal not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    constraint "order_items_pkey" primary key ("id")
);

create table "payments" (
    "id" bigserial not null,
    "payment_method" varchar not null,
    "image_url" varchar null,
    "amount" decimal not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    constraint "payments_pkey" primary key ("id")
);

alter table "orders"
add foreign key ("user_id") references "users" ("id");

alter table "orders"
add foreign key ("pharmacy_id") references "pharmacies" ("id");

alter table "orders"
add foreign key ("payment_id") references "payments" ("id");

alter table "orders"
add foreign key ("logistic_id") references "logistics" ("id");

alter table "order_items"
add foreign key ("order_id") references "orders" ("id");

alter table "order_items"
add foreign key ("product_id") references "products" ("id");