-- migrate:up
create schema if not exists addons;

create extension if not exists btree_gin schema addons;
create extension if not exists btree_gist schema addons;
-- https://www.postgresql.org/docs/current/ltree.html
create extension if not exists ltree schema addons;
-- https://www.postgresql.org/docs/current/pgtrgm.html
create extension if not exists pg_trgm schema addons;
-- https://www.postgresql.org/docs/current/pgcrypto.html
create extension if not exists pgcrypto schema addons;
-- https://www.postgresql.org/docs/current/unaccent.html
create extension if not exists unaccent schema addons;
-- create extension pg_cron schema addons;
-- https://docs.aws.amazon.com/amazonrds/latest/userguide/postgresql_pg_cron.html
-- grant usage on schema cron to postgres;

-- migrate:down