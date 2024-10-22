alter table if exists "password_reset_tokens" drop constraint "password_reset_tokens_user_id_fkey";

alter table if exists "verification_tokens" drop constraint "verification_tokens_user_id_fkey";

drop table if exists "password_reset_tokens";

drop table if exists "verification_tokens";