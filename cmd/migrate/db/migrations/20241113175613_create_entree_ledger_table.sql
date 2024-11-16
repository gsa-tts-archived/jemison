-- migrate:up

-- There is no 'if not exists' for CREATE TYPE
-- https://stackoverflow.com/questions/7624919/check-if-a-user-defined-type-already-exists-in-postgresql
-- Also, for future migrations:
-- https://stackoverflow.com/questions/1771543/adding-a-new-value-to-an-existing-enum-type

DO $$ BEGIN
  CREATE TYPE scheme AS ENUM (
    'http',
    'https'
  );
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS hosts (
  id BIGINT generated always as identity primary key,
  host TEXT,
  UNIQUE(host)
);

CREATE TABLE IF NOT EXISTS content_types (
  id INTEGER generated always as identity primary key,
  content_type TEXT
);

CREATE TABLE IF NOT EXISTS guestbook (
  id BIGINT generated always as identity primary key,
  scheme scheme NOT NULL,
  host BIGINT references hosts(id) NOT NULL,
  path TEXT NOT NULL,
  content_sha1 TEXT,
  content_length INTEGER,
  content_type INTEGER references content_types(id),
  last_updated DATE,
  last_fetched DATE,
  next_fetch DATE,
  UNIQUE (host, path)
);

-- migrate:down
DROP TABLE IF EXISTS guestbook;
DROP TABLE IF EXISTS hosts;
DROP TABLE IF EXISTS content_types;
DROP TYPE IF EXISTS scheme;


