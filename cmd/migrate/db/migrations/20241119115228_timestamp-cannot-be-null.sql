-- migrate:up
ALTER TABLE hosts ALTER COLUMN next_fetch SET NOT NULL;

-- migrate:down

