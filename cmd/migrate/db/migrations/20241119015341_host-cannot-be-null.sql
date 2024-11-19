-- migrate:up
ALTER TABLE hosts ALTER COLUMN host SET NOT NULL;

-- migrate:down

