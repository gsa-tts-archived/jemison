-- migrate:up
ALTER TABLE guestbook 
  ALTER COLUMN last_updated 
  TYPE TIMESTAMP;

ALTER TABLE guestbook 
  ALTER COLUMN last_fetched
  TYPE TIMESTAMP;

-- migrate:down
