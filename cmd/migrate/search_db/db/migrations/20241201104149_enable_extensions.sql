-- migrate:up
create schema if not exists addons;

CREATE EXTENSION btree_gin SCHEMA addons;
CREATE EXTENSION btree_gist SCHEMA addons;
-- https://www.postgresql.org/docs/current/ltree.html
CREATE EXTENSION ltree SCHEMA addons;
-- https://www.postgresql.org/docs/current/pgtrgm.html
CREATE EXTENSION pg_trgm SCHEMA addons;
-- https://www.postgresql.org/docs/current/pgcrypto.html
CREATE EXTENSION pgcrypto SCHEMA addons;
-- https://www.postgresql.org/docs/current/unaccent.html
CREATE EXTENSION unaccent SCHEMA addons;
-- CREATE EXTENSION pg_cron SCHEMA addons;
-- https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/PostgreSQL_pg_cron.html
-- GRANT USAGE ON SCHEMA cron TO postgres;

-- migrate:down