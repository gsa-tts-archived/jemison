-- migrate:up
CREATE TABLE if not exists metadata (
  id BIGSERIAL PRIMARY KEY,
  col TEXT,
  value TEXT
)
;

INSERT INTO metadata 
  (col, value)
  VALUES
  ('version', '1'),
  ('datetime_created', NOW()::TEXT),
  ('date_created', NOW()::DATE)
;  

CREATE TABLE if not exists paths (
  id BIGSERIAL PRIMARY KEY,
  host TEXT NOT NULL,
  path TEXT NOT NULL,
  UNIQUE(host, path)
);

CREATE TABLE  if not exists titles (
  id BIGSERIAL PRIMARY KEY,
  path_id INTEGER NOT NULL REFERENCES paths(id),
  kind INTEGER NOT NULL,
  title TEXT NOT NULL
);

CREATE TABLE  if not exists headers (
  id BIGSERIAL PRIMARY KEY,
  path_id INTEGER NOT NULL REFERENCES paths(id),
  kind INTEGER NOT NULL,
  level INTEGER NOT NULL,
  header TEXT NOT NULL
);

CREATE TABLE if not exists bodies (
  id BIGSERIAL PRIMARY KEY,
  path_id INTEGER NOT NULL REFERENCES paths(id),
  kind INTEGER NOT NULL,
  tag TEXT NOT NULL,
  body TEXT NOT NULL
);

-- All indicies
CREATE INDEX if not exists idx_gin_paths ON paths USING GIN (to_tsvector('english', path));
CREATE INDEX if not exists idx_gist_paths ON paths USING GiST (to_tsvector('english', path));

CREATE INDEX if not exists idx_gin_titles ON titles USING GIN (to_tsvector('english', title));
CREATE INDEX if not exists idx_gist_titles ON titles USING GiST (to_tsvector('english', title));

CREATE INDEX if not exists idx_gin_headers ON headers USING GIN (to_tsvector('english', header));
CREATE INDEX if not exists idx_gist_headers ON headers USING GiST (to_tsvector('english', header));

CREATE INDEX if not exists idx_gin_bodies ON bodies USING GIN (to_tsvector('english', body));
CREATE INDEX if not exists idx_gist_bodies ON bodies USING GiST (to_tsvector('english', body));

-- migrate:down

