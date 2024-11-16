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
  ON CONFLICT (host) DO NOTHING
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

-- If it is time to fetch, 
-- name: ToFetchOrNotToFetch :one
SELECT 
  CASE
    -- If we get NULL for a host/path, then it is time to fetch.
    -- It must be a page never seen before
    WHEN host = $1 and path = $2 and next_fetch = NULL THEN 'FETCH'
    -- We fetch after our time has past. So, it is not yet time to fetch.
    -- Only fetch if the date is STRICTLY in the past.
    WHEN host = $1 and path = $2 and next_fetch < NOW() THEN 'FETCH'
    -- If next_fetch is in our future (or today), do not fetch.
    WHEN host = $1 and path = $2 and next_fetch >= NOW() THEN 'DO_NOT_FETCH'
    ELSE 'FETCH_NOT_FOUND'
  END AS is_past_next_fetch
FROM guestbook; 

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