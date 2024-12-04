-- As we are building an index over a page, the first thing
-- that needs to happen is we insert a unique reference
-- in the path table.
-- name: InsertPath :one
INSERT INTO paths (path)
    SELECT $1
    WHERE NOT EXISTS (
        SELECT id, path
        FROM paths
        WHERE path = $1
    )
    RETURNING id
;

-- INSERT INTO tag ("key", "value")
-- SELECT 'key1', 'value1'
-- WHERE NOT EXISTS (
--     SELECT id, "key", "value"
--     FROM node_tag
--     WHERE key = 'key1' AND value = 'value1'
--     )
-- returning id, "key", "value"

-- name: InsertTitle :one
INSERT INTO titles
    (path_id, kind, title)
    VALUES ($1, 0, $2)
    RETURNING id
;

-- name: InsertHeader :one
INSERT INTO headers
    (path_id, kind, level, header)
    VALUES ($1, 1, $2, $3)
    RETURNING id
;

-- name: InsertBody :one
INSERT INTO bodies
    (path_id, kind, tag, body)
    VALUES
    ($1, 2, $2, $3)
    RETURNING id
;

-- name: CountSiteIndex :one
SELECT count(*) FROM paths;
