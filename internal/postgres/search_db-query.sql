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

-- name: SearchContent :many
select 
  to_hex(domain64), 
  path, 
  ts_rank_cd(contentsearch, query) AS Rank,
  ts_headline('english', sc.content, query, 'MaxFragments=10, MaxWords=7, MinWords=3') as Snippet,
  (select sc.content from searchable_content sc where sc.domain64 = @d64_start::bigint and sc.tag = 'title' limit 1) as Title
from searchable_content sc, 
	to_tsquery('english', @query::text) query,
	to_tsvector('english', sc.content) contentsearch
where query @@ contentsearch
and
domain64 >= @d64_start::bigint 
and domain64 < @d64_end::bigint
order by Rank desc
limit 10
;