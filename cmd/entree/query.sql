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

-- Does what it says; returns boolean. 
-- name: ItIsNotYetTimeToFetch :one
SELECT 
  CASE
    -- We fetch after our time has past. So, it is not yet time to fetch.
    WHEN host = $1 and next_fetch <= NOW() THEN 0::BOOLEAN
    ELSE 1::BOOLEAN
  END AS is_past_next_fetch
FROM guestbook; 