-- migrate:up

create table searchable_content (
  -- id bigint generate identity 
  domain64 bigint not null,
  path text not null,
  tag text default 'p' not null,
  content text not null
  -- unique(domain64)
)
-- partition by range(domain64)
;

-- https://www.postgresql.org/docs/current/pgtrgm.html#PGTRGM-TEXT-SEARCH
-- NOTE: Since the words table has been generated as a separate, static table, 
-- it will need to be periodically regenerated so that it remains reasonably 
-- up-to-date with the document collection. Keeping it exactly current is 
-- usually unnecessary.
-- create table words as select word from
--         ts_stat('select to_tsvector(''simple'', content) from searchable_content');

-- create index words_idx on words using gin (word gin_trgm_ops);

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

