#!/bin/bash

QUEUE_URL="postgresql://postgres@localhost:5432/postgres?sslmode=disable"

for i in $(seq 1 20);
do
  psql "${QUEUE_URL}" -c "truncate table river_job; truncate table river_leader; truncate table river_queue; select count(*) from river_job;"
done