--------------------------------------------------------------
-- `guestbook` table
--------------------------------------------------------------

  -- id BIGINT generated always as identity primary key,
  -- scheme scheme NOT NULL,
  -- host BIGINT references hosts(id) NOT NULL,
  -- path TEXT NOT NULL,
  -- content_sha1 TEXT,
  -- content_length INTEGER,
  -- content_type INTEGER references content_types(id),
  -- last_updated TIMESTAMP NOT NULL,
  -- last_fetched TIMESTAMP NOT NULL,
  -- next_fetch TIMESTAMP NOT NULL,
  -- UNIQUE (host, path)

-- Returns a single guestbook entry
-- based on the unique host/path combo.
-- name: GetGuestbookEntry :one
SELECT * FROM guestbook 
WHERE 
  host = $1 
  AND
  path = $2
  LIMIT 1;

-- Returns all the entries for a host
-- name: GetGuestbookEntries :many
SELECT path FROM guestbook 
WHERE 
  host = $1;

-- name: CheckEntryExistsInGuestbook :one
SELECT EXISTS(SELECT 1::bool FROM guestbook WHERE host = $1);

-- This is likely called by `fetch`
-- name: UpdateGuestbookFetch :one
INSERT INTO guestbook
  -- 1       2     3        4             -           5
  (scheme, host, path, last_updated, last_fetched, next_fetch)
  VALUES 
  -- always insert three seconds in the past
  -- this way, if we check too fast, we're guaranteed to see
  -- that we last fetched it in the past.
  ($1, $2, $3, $4, NOW()+make_interval(secs => -3), $5)
  ON CONFLICT (host, path) DO UPDATE
  SET
    scheme = EXCLUDED.scheme,
    host = EXCLUDED.host,
    path = EXCLUDED.path,
    content_sha1 = EXCLUDED.content_sha1,
    content_length = EXCLUDED.content_length,
    content_type = EXCLUDED.content_type,
    last_updated = EXCLUDED.last_updated,
    last_fetched = EXCLUDED.last_fetched,
    next_fetch = EXCLUDED.next_fetch
  RETURNING id;



--------------------------------------------------------------
-- `host` table
--------------------------------------------------------------
-- Find a host ID
-- name: GetHostId :one
SELECT id FROM hosts WHERE host = $1;

-- name: GetHostNextFetch :one
SELECT next_fetch from hosts WHERE id = $1;


-- This gets used at startup. We walk the config
-- file and load the table with unique hosts, so that
-- we can use the IDs in the `guestbook` table.
-- name: UpsertUniqueHost :one
INSERT INTO hosts 
  (host, next_fetch) 
  VALUES ($1, $2) 
  -- See https://stackoverflow.com/a/37543015
  -- This is a workaround; it forces the `id` to be 
  -- returned in all cases.
  ON CONFLICT (host) DO UPDATE
    SET host = EXCLUDED.host,
    next_fetch = EXCLUDED.next_fetch
  RETURNING id
;

  -- (SELECT
  --   CASE 
  --     WHEN $1::TEXT = 'weekly' THEN NOW()+make_interval(weeks => 1, days => -1)
  --     WHEN $1 = 'bi-weekly' THEN NOW()+make_interval(weeks => 2, days => -1)
  --     WHEN $1 = 'monthly' THEN NOW()+make_interval(months => 1, days => -1)
  --     WHEN $1 = 'bi-monthly' THEN NOW()+make_interval(months => 2, days => -1)
  --     WHEN $1 = 'quarterly' THEN NOW()+make_interval(months=> 3, days => -1)
  --   END)
  -- )

-- FIXM: This is the same as Upsert. Reduce to one function.
-- name: UpdateHostNextFetch :one
INSERT INTO hosts
  (host, next_fetch)
  VALUES ($1, $2)
  ON CONFLICT (host)
  WHERE host = $1
  DO UPDATE
    SET 
      host = EXCLUDED.host,
      next_fetch = EXCLUDED.next_fetch
  RETURNING id
;

-- name: SetHostNextFetchToYesterday :exec
UPDATE hosts
  SET
    next_fetch = NOW()+make_interval(days => -1)
  WHERE
    host = $1
;

-- name: SetGuestbookFetchToYesterdayForHost :exec
UPDATE guestbook
  SET 
    next_fetch = NOW()+make_interval(days => -1)
  WHERE
    host = $1
;

-- If it is time to fetch?
-- These match a Golang enum.
-- const (
-- 	NotPreviouslySeen FetchStatus = iota 0
-- 	DeadlineNotYetReached = 1
-- 	DeadlinePassed = 2
-- 	NotFound = 3
-- 	HallPass = 4 
-- 	DefaultCase = 5
-- )
-- name: ToFetchOrNotToFetch :one
-- SELECT 
--   CASE
--     -- Deadline passed
--     WHEN host = $1 AND path = $2 AND next_fetch < NOW() THEN 'DeadlinePassed'
--     -- Deadline not yet reached
--     WHEN host = $1 AND path = $2 AND next_fetch >= NOW() THEN 'DeadlineNotYetReached'
--     -- If we can't find the path, then this is not found
--     ELSE 'NotFound'
--   END AS is_past_next_fetch
--   FROM guestbook
--   WHERE host = $1 AND path = $2
-- ;

