--------------------------------------------------------------
-- `guestbook` table
--------------------------------------------------------------

-- returns a single guestbook entry
-- based on the unique host/path combo.
-- name: GetGuestbookEntry :one
select * from guestbook 
where 
  domain64 = $1 
  and
  path = $2
  limit 1;

-- returns all the entries for a host
-- name: GetGuestbookEntries :many
select path from guestbook 
where 
  domain64 = $1;

-- name: CheckEntryExistsInGuestbook :one
select exists(select 1::bool from guestbook where domain64 = $1);

-- this is likely called by `entree`
-- name: UpdateGuestbookNextFetch :one
insert into guestbook
  -- 1       2         3        4
  (scheme, domain64, path, next_fetch)
  values 
  ($1, $2, $3, $4)
  on conflict (domain64, path) do update
  set
    scheme = $1,
    domain64 = $2,
    path = $3,
    content_length = excluded.content_length,
    content_type = excluded.content_type,
    last_modified = excluded.last_modified,
    last_fetched = now()+make_interval(secs => -3),
    next_fetch = coalesce($4, excluded.next_fetch)
  returning id
;

-- this is likely called by `fetch`
-- name: UpdateGuestbookFetch :one
insert into guestbook
  -- 1       2          3        4
  (scheme, domain64, path, content_length, 
  --    5               6             7         8
  content_type, last_modified, last_fetched, next_fetch)
  values 
  -- always insert three seconds in the past
  -- this way, if we check too fast, we're guaranteed to see
  -- that we last fetched it in the past.
  ($1, $2, $3, $4, $5, $6, $7, $8)
  on conflict (domain64, path) do update
  set
    scheme = excluded.scheme,
    domain64 = excluded.domain64,
    path = excluded.path,
    content_length = coalesce($4, excluded.content_length),
    content_type = coalesce($5, excluded.content_type),
    last_modified = coalesce($6, excluded.last_modified),
    -- last_fetched = coalesce($7, now()+make_interval(secs => -3)),
    last_fetched = coalesce($7, excluded.last_fetched),
    next_fetch = coalesce($8, excluded.next_fetch)
  returning id
;

-- --------------------------------------------------------------
-- -- `host` table
-- --------------------------------------------------------------
-- -- find a host id
-- -- name: GetHostId :one
-- select id from hosts where host = $1;

-- name: GetHostNextFetch :one
select next_fetch from hosts where domain64 = $1::bigint;


-- name: SetHostNextFetchToYesterday :exec
update hosts
  set
    next_fetch = now()+make_interval(days => -1)
  where
    domain64 = $1::bigint
;

-- name: SetGuestbookFetchToYesterdayForHost :exec
update guestbook
  set 
    next_fetch = now()+make_interval(days => -1)
  where
    domain64 = $1::bigint
;


-- this gets used at startup. we walk the config
-- file and load the table with unique hosts, so that
-- we can use the ids in the `guestbook` table.
-- name: UpsertUniqueHost :one
insert into hosts 
  (domain64, next_fetch) 
  values ($1, $2) 
  -- see https://stackoverflow.com/a/37543015
  -- this is a workaround; it forces the `id` to be 
  -- returned in all cases.
  on conflict (domain64) do update
    set domain64 = coalesce($1::bigint, excluded.domain64),
    next_fetch = excluded.next_fetch
  returning id
;

-- name: PathsInDomain64Range :one
select count(*)
from guestbook
where
  domain64 >= @d64_start::bigint 
  and 
  domain64 < @d64_end::bigint
;