create table "customer_details" (
    "user_id" bigint not null,
    "profile_photo_url" varchar null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    constraint "customer_details_pkey" primary key ("user_id")
);

create table "pharmacist_details" (
    "user_id" bigint not null,
    "name" varchar not null,
    "sipa_number" varchar unique not null,
    "whatsapp_number" varchar unique not null,
    "years_of_experience" varchar not null,
    "is_assigned" boolean not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    constraint "pharmacist_details_pkey" primary key ("user_id")
);

alter table "customer_details"
add foreign key ("user_id") references "users" ("id");

alter table "pharmacist_details"
add foreign key ("user_id") references "users" ("id");