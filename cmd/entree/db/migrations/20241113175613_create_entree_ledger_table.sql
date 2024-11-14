-- migrate:up
CREATE TYPE scheme AS ENUM (
  'http',
  'https'
);

-- CREATE TABLE IF NOT EXISTS schemes (
--   id INTEGER generated always as identity primary key,
--   scheme enum_scheme
-- );

CREATE TABLE IF NOT EXISTS hosts (
  id BIGINT generated always as identity primary key,
  host TEXT
);

CREATE TABLE IF NOT EXISTS content_types (
  id INTEGER generated always as identity primary key,
  content_type TEXT
);

CREATE TABLE IF NOT EXISTS ledger (
  id BIGINT generated always as identity primary key,
  scheme scheme NOT NULL,
  host BIGINT references hosts(id),
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

