#!/bin/bash

GBURL="postgresql://postgres@localhost:6543/postgres?sslmode=disable"
QUEUE_URL="postgresql://postgres@localhost:5432/postgres?sslmode=disable"

psql "${QUEUE_URL}" -c "select queue, state, count(*) from river_job where state='available' group by queue, state order by queue"

