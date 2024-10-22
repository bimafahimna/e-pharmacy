create table "products" (
    "id" bigserial not null,
    "manufacturer_id" bigint not null,
    "product_classification_id" bigint not null,
    "product_form_id" bigint null,
    "name" varchar not null,
    "generic_name" varchar not null,
    "categories" varchar[] null,
    "description" text not null,
    "unit_in_pack" int null,
    "selling_unit" varchar null,
    "daily_sold_amount" int not null default 0,
    "sold_amount" int not null default 0,
    "weight" decimal not null,
    "height" decimal not null,
    "length" decimal not null,
    "width" decimal not null,
    "image_url" varchar not null,
    "usage" int not null,
    "is_active" boolean null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    "deleted_at" timestamp null,
    constraint "products_pkey" primary key ("id")   
);

create table "categories" (
    "id" bigserial not null,
    "name" varchar unique not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp
);

create table "manufacturers" (
    "id" bigserial not null,
    "name" varchar not null,
    constraint "manufacturers_pkey" primary key ("id")
);

create table "product_classifications" (
    "id" bigserial not null,
    "name" varchar not null,
    constraint "product_classifications_pkey" primary key ("id")
);

create table "product_forms" (
    "id" bigserial not null,
    "name" varchar not null,
    constraint "product_forms_pkey" primary key ("id")
);

alter table "products"
add foreign key ("manufacturer_id") references "manufacturers" ("id");

alter table "products"
add foreign key ("product_classification_id") references "product_classifications" ("id");

alter table "products"
add foreign key ("product_form_id") references "product_forms" ("id");

alter table "products"
add constraint "products_name_generic_name_manufacturer_id_key" unique ("name", "generic_name", "manufacturer_id");

create index products_index_sort_by on products ("id","daily_sold_amount" desc nulls last, "sold_amount" desc nulls last) where "is_active" = true;
