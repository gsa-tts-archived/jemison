-- migrate:up

-- There is no 'if not exists' for CREATE TYPE
-- https://stackoverflow.com/questions/7624919/check-if-a-user-defined-type-already-exists-in-postgresql
-- Also, for future migrations:
-- https://stackoverflow.com/questions/1771543/adding-a-new-value-to-an-existing-enum-type

CREATE TABLE IF NOT EXISTS hosts (
  id BIGINT generated always as identity primary key,
  host TEXT NOT NULL,
  next_fetch TIMESTAMP NOT NULL,
  UNIQUE(host)
);

-- CREATE TABLE IF NOT EXISTS content_types (
--   id INTEGER generated always as identity primary key,
--   content_type TEXT
--   UNIQUE(content_type)
-- );

CREATE TABLE IF NOT EXISTS guestbook (
  id BIGINT generated always as identity primary key,
  scheme TEXT NOT NULL DEFAULT 'https',
  host BIGINT references hosts(id) NOT NULL,
  path TEXT NOT NULL,
  content_sha1 TEXT,
  content_length INTEGER,
  content_type TEXT,
  last_updated TIMESTAMP NOT NULL,
  last_fetched TIMESTAMP NOT NULL,
  next_fetch TIMESTAMP NOT NULL,
  -- FIXME: Add ETag
  UNIQUE (host, path)
);

-- migrate:down
DROP TABLE IF EXISTS guestbook;
DROP TABLE IF EXISTS hosts;
DROP TABLE IF EXISTS content_types;


