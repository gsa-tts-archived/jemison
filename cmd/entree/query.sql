-- -- name: CountJobs :one
-- SELECT count(*) from river_job WHERE kind=$1 AND state=$2;

  -- id integer generated always as identity primary key,
  -- scheme INTEGER references schemes(id),
  -- host INTEGER references hosts(id),
  -- path TEXT,
  -- sha1 TEXT,
  -- length INTEGER,
  -- content_type TEXT,
  -- last_updated DATE,
  -- last_fetched DATE,
  -- next_fetch DATE

-- This is so we can load all of the hosts into an 
-- in-memory dictionary at startup. There will only be
-- a few thousand, and therefore is easy to load and
-- then use to lookup host IDs.
-- name: GetHosts :many
SELECT * from hosts;

-- name: GetLedger :one
SELECT * FROM ledger 
WHERE 
  host = $1 
  AND
  path = $2
  LIMIT 1;
