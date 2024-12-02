#!/bin/bash

GBURL="postgresql://postgres@localhost:6543/postgres?sslmode=disable"

for i in $(seq 1 3);
do
  psql "${GBURL}" -c "truncate table guestbook"
done

psql "${GBURL}" -c "drop table guestbook ; drop table hosts ; truncate table schema_migrations;"
