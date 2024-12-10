--------------------------------------------------------------
-- `guestbook` table
--------------------------------------------------------------

  -- id bigint generated always as identity primary key,
  -- scheme scheme not null,
  -- host bigint references hosts(id) not null,
  -- path text not null,
  -- content_sha1 text,
  -- content_length integer,
  -- content_type integer references content_types(id),
  -- last_updated timestamp not null,
  -- last_fetched timestamp not null,
  -- next_fetch timestamp not null,
  -- unique (host, path)

-- returns a single guestbook entry
-- based on the unique host/path combo.
-- name: GetGuestbookEntry :one
select * from guestbook 
where 
  host = $1 
  and
  path = $2
  limit 1;

-- returns all the entries for a host
-- name: GetGuestbookEntries :many
select path from guestbook 
where 
  host = $1;

-- name: CheckEntryExistsInGuestbook :one
select exists(select 1::bool from guestbook where host = $1);

-- this is likely called by `entree`
-- name: UpdateGuestbookNextFetch :one
insert into guestbook
  -- 1       2     3        4
  (scheme, host, path, next_fetch)
  values 
  ($1, $2, $3, $4)
  on conflict (host, path) do update
  set
    scheme = excluded.scheme,
    host = excluded.host,
    path = excluded.path,
    content_length = coalesce($4, excluded.content_length),
    content_type = coalesce($5, excluded.content_type),
    last_modified = coalesce(excluded.last_updated),
    last_fetched = now()+make_interval(secs => -3),
    next_fetch = excluded.next_fetch
  returning id
;

-- this is likely called by `fetch`
-- name: UpdateGuestbookFetch :one
insert into guestbook
  -- 1       2     3        4
  (scheme, host, path, content_length, 
  --    5               6             7         8
  content_type, last_modified, last_fetched, next_fetch)
  values 
  -- always insert three seconds in the past
  -- this way, if we check too fast, we're guaranteed to see
  -- that we last fetched it in the past.
  ($1, $2, $3, $4, $5, $6, $7, $8)
  on conflict (host, path) do update
  set
    scheme = excluded.scheme,
    host = excluded.host,
    path = excluded.path,
    content_length = coalesce($4, excluded.content_length),
    content_type = coalesce($5, excluded.content_type),
    last_modified = coalesce(excluded.last_updated),
    last_fetched = now()+make_interval(secs => -3),
    next_fetch = excluded.next_fetch
  returning id
;

----------------------------------------
-- constants
----------------------------------------

-- Return `https` as the default
-- name: GetScheme :one
select coalesce(id, 1) 
  from schemes 
  where scheme = $1
;

-- Return `gov` as the default
-- name: GetTLD :one
select coalesce(id, 1)
  from tlds
  where tld = $1
;

-- select `binary/octet-stream` 
-- name: GetContentType :one
select coalesce(id, 1)
  from content_types
  where content_type = $1
;

-- `p` is the default
-- name: GetTag :one
select coalesce(id, 1)
  from tags
  where tag = $1
;






-- --------------------------------------------------------------
-- -- `host` table
-- --------------------------------------------------------------
-- -- find a host id
-- -- name: GetHostId :one
-- select id from hosts where host = $1;

-- -- name: GetHostNextFetch :one
-- select next_fetch from hosts where id = $1;


-- -- this gets used at startup. we walk the config
-- -- file and load the table with unique hosts, so that
-- -- we can use the ids in the `guestbook` table.
-- -- name: UpsertUniqueHost :one
-- insert into hosts 
--   (host, next_fetch) 
--   values ($1, $2) 
--   -- see https://stackoverflow.com/a/37543015
--   -- this is a workaround; it forces the `id` to be 
--   -- returned in all cases.
--   on conflict (host) do update
--     set host = excluded.host,
--     next_fetch = excluded.next_fetch
--   returning id
-- ;

--   -- (select
--   --   case 
--   --     when $1::text = 'weekly' then now()+make_interval(weeks => 1, days => -1)
--   --     when $1 = 'bi-weekly' then now()+make_interval(weeks => 2, days => -1)
--   --     when $1 = 'monthly' then now()+make_interval(months => 1, days => -1)
--   --     when $1 = 'bi-monthly' then now()+make_interval(months => 2, days => -1)
--   --     when $1 = 'quarterly' then now()+make_interval(months=> 3, days => -1)
--   --   end)
--   -- )

-- -- fixm: this is the same as upsert. reduce to one function.
-- -- name: UpdateHostNextFetch :one
-- insert into hosts
--   (host, next_fetch)
--   values ($1, $2)
--   on conflict (host)
--   where host = $1
--   do update
--     set 
--       host = excluded.host,
--       next_fetch = excluded.next_fetch
--   returning id
-- ;

-- -- name: SetHostNextFetchToYesterday :exec
-- update hosts
--   set
--     next_fetch = now()+make_interval(days => -1)
--   where
--     host = $1
-- ;

-- -- name: SetGuestbookFetchToYesterdayForHost :exec
-- update guestbook
--   set 
--     next_fetch = now()+make_interval(days => -1)
--   where
--     host = $1
-- ;
