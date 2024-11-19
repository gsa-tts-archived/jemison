-- Find a host ID
-- name: GetHostId :one
SELECT id FROM hosts WHERE host = $1;

-- name: GetHostNextFetch :one
SELECT next_fetch from hosts WHERE id = $1;

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


-- This gets used at startup. We walk the config
-- file and load the table with unique hosts, so that
-- we can use the IDs in the `guestbook` table.
-- name: UpsertUniqueHost :one
INSERT INTO hosts 
  (host, next_fetch) 
  VALUES ($2, 
  (SELECT
    CASE 
      WHEN $1 = 'weekly' THEN NOW()+make_interval(weeks => 1, days => -1)
      WHEN $1 = 'bi-weekly' THEN NOW()+make_interval(weeks => 2, days => -1)
      WHEN $1 = 'monthly' THEN NOW()+make_interval(months => 1, days => -1)
      WHEN $1 = 'bi-monthly' THEN NOW()+make_interval(months => 2, days => -1)
      WHEN $1 = 'quarterly' THEN NOW()+make_interval(months=> 3, days => -1)
    END)
  )
  -- See https://stackoverflow.com/a/37543015
  -- This is a workaround; it forces the `id` to be 
  -- returned in all cases.
  ON CONFLICT (host) DO UPDATE
    SET host = EXCLUDED.host,
    next_fetch = EXCLUDED.next_fetch
  RETURNING id
;

-- name: UpdateNextFetch :one
INSERT INTO hosts
  (host, next_fetch)
  VALUES ($2, 
  (SELECT
    CASE 
      WHEN $1 = 'weekly' THEN NOW()+make_interval(weeks => 1, days => -1)
      WHEN $1 = 'bi-weekly' THEN NOW()+make_interval(weeks => 2, days => -1)
      WHEN $1 = 'monthly' THEN NOW()+make_interval(months => 1, days => -1)
      WHEN $1 = 'bi-monthly' THEN NOW()+make_interval(months => 2, days => -1)
      WHEN $1 = 'quarterly' THEN NOW()+make_interval(months=> 3, days => -1)
    END)
  )
  ON CONFLICT (host)
  WHERE host = $2
  DO UPDATE
    SET 
      host = EXCLUDED.host,
      next_fetch = EXCLUDED.next_fetch
  RETURNING id
;

-- name: SetNextFetchToYesterday :one
INSERT INTO hosts
  (host, next_fetch)
  VALUES ($1, NOW()+make_interval(days => -1))
  ON CONFLICT (host)
  WHERE host = $1
  DO UPDATE
    SET 
      host = EXCLUDED.host,
      next_fetch = EXCLUDED.next_fetch
  RETURNING id
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

