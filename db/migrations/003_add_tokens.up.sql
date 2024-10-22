create table "verification_tokens" (
    "id" bigserial not null,
    "user_id" bigint not null,
    "token" varchar unique not null,
    "used_at" timestamp null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    "expired_at" timestamp not null,
    constraint "verification_tokens_pkey" primary key ("id")
);

create table "password_reset_tokens" (
    "id" bigserial not null,
    "user_id" bigint not null,
    "token" varchar unique not null,
    "used_at" timestamp null,
    "created_at" timestamp not null default current_timestamp,
    "updated_at" timestamp not null default current_timestamp,
    "expired_at" timestamp not null,
    constraint "password_reset_tokens_pkey" primary key ("id")
);

alter table "verification_tokens"
add foreign key ("user_id") references "users" ("id");

alter table "password_reset_tokens"
add foreign key ("user_id") references "users" ("id");