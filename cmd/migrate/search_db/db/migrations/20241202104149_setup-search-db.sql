-- migrate:up
create table if not exists metadata (
  id integer generated always as identity primary key,
  col text,
  value text
)
;

insert into metadata 
  (col, value)
  values
  ('version', '1'),
  ('datetime_created', now()::text),
  ('date_created', now()::date)
;  

create table raw_content (
  id bigint generated always as identity primary key,
  domain64 bigint not null,
  path text not null,
  tag text default 'p' not null,
  content text not null
)
;

-- create table if not exists paths (
--   id bigserial primary key,
--   host text not null,
--   path text not null,
--   unique(host, path)
-- );

-- create table  if not exists titles (
--   id bigserial primary key,
--   path_id integer not null references paths(id),
--   kind integer not null,
--   title text not null
-- );

-- create table  if not exists headers (
--   id bigserial primary key,
--   path_id integer not null references paths(id),
--   kind integer not null,
--   level integer not null,
--   header text not null
-- );

-- create table if not exists bodies (
--   id bigserial primary key,
--   path_id integer not null references paths(id),
--   kind integer not null,
--   tag text not null,
--   body text not null
-- );

-- -- all indicies
-- create index if not exists idx_gin_paths on paths using gin (to_tsvector('english', path));
-- create index if not exists idx_gist_paths on paths using gist (to_tsvector('english', path));

-- create index if not exists idx_gin_titles on titles using gin (to_tsvector('english', title));
-- create index if not exists idx_gist_titles on titles using gist (to_tsvector('english', title));

-- create index if not exists idx_gin_headers on headers using gin (to_tsvector('english', header));
-- create index if not exists idx_gist_headers on headers using gist (to_tsvector('english', header));

-- create index if not exists idx_gin_bodies on bodies using gin (to_tsvector('english', body));
-- create index if not exists idx_gist_bodies on bodies using gist (to_tsvector('english', body));

-- migrate:down

