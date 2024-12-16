-- name: getrawcontent :one
select * from raw_content 
where 
  domain64 = $1 
  and
  path = $2
  limit 1;


---------------------------------------------------------------
-- raw_content
---------------------------------------------------------------
-- Create table raw_content (
--   id bigint generated always as identity primary key,
--   domain64 bigint,
--   path text,
--   tag text default 'p',
--   content text
-- )
-- ;

-- name: InsertRawContent :exec
insert into raw_content
  (domain64, path, tag, content)
  values
  ($1, $2, $3, $4)
;

-- select 
--   to_hex(domain64), 
--   path, 
--   ts_rank_cd(contentsearch, query) AS Rank,
--   ts_headline('english', sc.content, query, 'MaxFragments=3, MaxWords=7, MinWords=3') as Snippet
-- from searchable_content sc, 
-- 	to_tsquery('english', @query::text) query,
-- 	to_tsvector('english', sc.content) contentsearch
-- where 
--   query @@ contentsearch
--   and
--   domain64 >= @d64_start::bigint 
--   and 
--   domain64 < @d64_end::bigint
-- order by Rank desc
-- limit 10
-- ;

-- name: SearchContent :many
select 
  to_hex(domain64), 
  path, 
  ts_rank_cd(fts, query) AS Rank,
  ts_headline('english', sc.content, query, 'MaxFragments=3, MaxWords=7, MinWords=3') as Snippet
from searchable_content sc, 
	to_tsquery('english', @query::text) query
where 
  domain64 >= @d64_start::bigint 
  and 
  domain64 < @d64_end::bigint
  and 
  fts @@ query
order by Rank desc
limit 10
;

-- name: GetTitle :one
select 
  sc.content
from searchable_content sc
where
  sc.domain64 = @domain64::bigint
  and
  sc.path = @path::text
  and
  sc.tag = 'title'
limit 1
;

-- name: GetPath :one
select 
  sc.content
from searchable_content sc
where
  sc.domain64 = @domain64::bigint
  and
  sc.path = @path::text
  and
  sc.tag = 'path'
limit 1
;

-- name: BodiesInDomain64Range :one
select count(*)
from searchable_content
where
  domain64 >= @d64_start::bigint 
  and 
  domain64 < @d64_end::bigint
  and
  tag = 'body'
;