-- migrate:up

create table searchable_content (
  domain64 bigint not null,
  path text not null,
  tag text default 'p' not null,
  content text not null
)
;

create index if not exists domain64_idx on searchable_content (domain64);
create index if not exists tag_idx on searchable_content (tag);

create index if not exists idx_gin_paths on searchable_content 
  using gin (to_tsvector('english', path));

create index if not exists idx_gist_paths on searchable_content 
  using gist (to_tsvector('english', path));

create index if not exists idx_gin_bodies on searchable_content 
  using gin (to_tsvector('english', content));

create index if not exists idx_gist_bodies on searchable_content 
  using gist (to_tsvector('english', content));

-- migrate:down

