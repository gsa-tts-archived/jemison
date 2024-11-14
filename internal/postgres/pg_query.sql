-- -- name: CountJobs :one
-- SELECT count(*) from river_job WHERE kind=$1 AND state=$2;

-- name: Update
