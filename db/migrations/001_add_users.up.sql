create table "users" (
    "id" bigserial not null,
    "role" varchar not null,
    "email" varchar unique not null,
    "password_hash" varchar null,
    "is_verified" boolean not null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    "deleted_at" timestamp null,
    constraint "users_pkey" primary key ("id")
);
