#!/bin/bash

GBURL="postgresql://postgres@localhost:6543/postgres?sslmode=disable"
QUEUE_URL="postgresql://postgres@localhost:5432/postgres?sslmode=disable"

for i in $(seq 1 3);
do
  psql "${GBURL}" -c "truncate table guestbook"
  psql "${QUEUE_URL}" -c "truncate table river_job; truncate table river_leader; truncate table river_queue; select count(*) from river_job;"
  sleep 0.3
  # psql "${GBURL}" -c "delete from hosts *"
  # psql "${GBURL}" -c "drop table guestbook ; drop table hosts ; truncate table schema_migrations;"
done


# for i in $(seq 1 3);
# do
#   psql "${QUEUE_URL}" -c "truncate table river_job; truncate table river_leader; truncate table river_queue; select count(*) from river_job;"
# done

# Reset minio

mc alias set localm http://localhost:9100 numbernine numbernine

for i in $(seq 1 3);
do
  mc rm --force --recursive localm/fetch
  mc rm --force --recursive localm/extract
  mc rm --force --recursive localm/serve
done