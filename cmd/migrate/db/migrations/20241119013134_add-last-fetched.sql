-- migrate:up
ALTER TABLE hosts
ADD COLUMN next_fetch TIMESTAMP NOT NULL
DEFAULT NOW()+make_interval(weeks => 1, days => -1);

ALTER TABLE guestbook
DROP COLUMN next_fetch;

-- migrate:down

ALTER TABLE hosts
DROP COLUMN next_fetch;

-- If we run this backwards, we lose the 
-- real timestamps. :shrug:
ALTER TABLE guestbook
ADD COLUMN next_fetch TIMESTAMP NOT NULL
DEFAULT NOW()+make_interval(weeks => 1, days => -1);
