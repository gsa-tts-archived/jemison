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


-- name: GetLedger :one
SELECT * FROM guestbook 
WHERE 
  host = $1 
  AND
  path = $2
  LIMIT 1;

-- name: UpsertUniqueHost :one
INSERT INTO hosts 
  (host) 
  VALUES ($1)
  ON CONFLICT DO NOTHING
  RETURNING id
;
