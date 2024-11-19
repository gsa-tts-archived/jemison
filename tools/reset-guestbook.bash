#!/bin/bash

GBURL="postgresql://postgres@localhost:6543/postgres?sslmode=disable"

for i in $(seq 1 20);
do
  psql "${GBURL}" -c "drop table guestbook ; drop table hosts ; truncate table schema_migrations;"
done

