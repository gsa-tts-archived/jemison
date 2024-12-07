-- name: getrawcontent :one
select * from raw_content 
where 
  domain64 = $1 
  and
  path = $2
  limit 1;