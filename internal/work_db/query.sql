-- Find a host ID
-- name: GetHostId :one
SELECT id FROM hosts WHERE host = $1;

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
SELECT path, next_fetch FROM guestbook 
WHERE 
  host = $1;


-- This gets used at startup. We walk the config
-- file and load the table with unique hosts, so that
-- we can use the IDs in the `guestbook` table.
-- name: UpsertUniqueHost :one
INSERT INTO hosts 
  (host) 
  VALUES ($1)
  -- See https://stackoverflow.com/a/37543015
  -- This is a workaround; it forces the `id` to be 
  -- returned in all cases.
  ON CONFLICT (host) DO UPDATE
    SET host = EXCLUDED.host
  RETURNING id
;

-- Update the next fetch time based on our schedule
-- name: UpdateNextFetch :one
INSERT INTO guestbook
  (scheme, host, path, last_fetched, next_fetch)
  VALUES ($2, $3, $4, NOW(),
  (SELECT
    CASE 
      WHEN $1 = 'weekly' THEN NOW()+make_interval(weeks => 1, days => -1)
      WHEN $1 = 'monthly' THEN NOW()+make_interval(months => 1, days => -1)
      WHEN $1 = 'bi-monthly' THEN NOW()+make_interval(months => 2, days => -1)
      WHEN $1 = 'quarterly' THEN NOW()+make_interval(months=> 3, days => -1)
    END)
  )
  ON CONFLICT (host, path) 
  WHERE host = $3 AND path = $4
  DO UPDATE 
    SET 
      next_fetch = EXCLUDED.next_fetch,
      last_fetched = NOW()
  RETURNING id;


-- If it is time to fetch?
-- These match a Golang enum.
-- const (
-- 	NotPreviouslySeen FetchStatus = iota
-- 	DeadlineNotYetReached
-- 	DeadlinePassed
-- 	NotFound
-- 	HallPass
-- 	DefaultCase
-- )
-- name: ToFetchOrNotToFetch :one
SELECT 
  CASE
    -- Deadline passed
    WHEN host = $1 AND path = $2 AND next_fetch < NOW() THEN 2
    -- Deadline not yet reached
    WHEN host = $1 AND path = $2 AND next_fetch >= NOW() THEN 1
    -- If we can't find the path, then this is not found
    ELSE 3
  END AS is_past_next_fetch
  FROM guestbook
  WHERE host = $1 AND path = $2
;


-- In order to crawl a site out of the normal schedule, we will want to 
-- rewrite when those pages were last fetched. We do this for a full domain.
-- name: SetLastFetchedToYesterday :one
INSERT INTO guestbook
  (host, path, next_fetch)
  VALUES ($1, $2, NOW()+make_interval(days => -1))
  ON CONFLICT (host, path)
  WHERE host = $1 AND path = $1
  DO UPDATE
    SET
      next_fetch = EXCLUDED.next_fetch
  RETURNING id;