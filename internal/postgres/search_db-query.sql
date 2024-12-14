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