#!/bin/bash

GBURL="postgresql://postgres@localhost:6543/postgres?sslmode=disable"

for i in $(seq 1 3);
do
  psql "${GBURL}" -c "truncate table guestbook"
  psql "${GBURL}" -c "delete from hosts *"
  # psql "${GBURL}" -c "drop table guestbook ; drop table hosts ; truncate table schema_migrations;"
done

QUEUE_URL="postgresql://postgres@localhost:5432/postgres?sslmode=disable"

for i in $(seq 1 3);
do
  psql "${QUEUE_URL}" -c "truncate table river_job; truncate table river_leader; truncate table river_queue; select count(*) from river_job;"
done

# Reset minio

mc alias set localm http://localhost:9000 numbernine numbernine
mc rm --force --recursive localm/fetch
mc rm --force --recursive localm/extract
mc rm --force --recursive localm/serve
