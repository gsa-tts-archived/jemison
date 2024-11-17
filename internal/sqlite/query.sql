-- As we are building an index over a page, the first thing
-- that needs to happen is we insert a unique reference
-- in the path table.
-- name: InsertPath :one
INSERT INTO paths 
    (path)
    VALUES (?)
    RETURNING id
;

-- name: InsertTitle :one
INSERT INTO titles
    (path_id, kind, title)
    VALUES (?, 0, ?)
    RETURNING id
;

-- name: InsertHeader :one
INSERT INTO headers
    (path_id, kind, level, header)
    VALUES (?, 1, ?, ?)
    RETURNING id
;

-- name: InsertBody :one
INSERT INTO bodies
    (path_id, kind, tag, body)
    VALUES
    (?, 2, ?, ?)
    RETURNING id
;

-- name: CountSiteIndex :one
SELECT count(*) FROM paths;

-- -- https://www.sqlitetutorial.net/sqlite-full-text-search/
-- -- name: SearchSiteIndexSnippets :many
-- SELECT path, snippet(site_index, 2, '<b>', '</b>', '...', 16)
--     FROM site_index 
--     WHERE text MATCH ?
--     AND path LIKE ?
--     ORDER BY rank
--     LIMIT ?;

-- name: Search :many
-- SELECT 
--   path_id,
--   (SELECT path from paths WHERE id = path_id), 
--   kind,
--   weight,
--   txt
--   FROM
--     (SELECT titles_fts.path_id as path_id, 4.0 as weight, kind, title as txt
--       FROM titles_fts
--       WHERE title MATCH ?1
--     UNION ALL
--     SELECT headers_fts.path_id as path_id, 2.0 as weight, kind, header as txt
--       FROM headers_fts
--       WHERE header MATCH ?1
--     UNION ALL
--     SELECT bodies_fts.path_id as path_id, 1.0 as weight, kind, body as txt
--       FROM bodies_fts
--       WHERE body MATCH ?1
--     ORDER BY weight DESC)

-- -- name: CreateSiteEntry :one
-- INSERT INTO site_index (host, path, text) 
--     VALUES (?, ?, ?)
--     RETURNING *;

-- -- This could become a AS (SELECT ...) 
-- -- so the table has the data at creation time?
-- -- name: SetVersion :exec
-- INSERT INTO metadata (version, last_updated) 
--     VALUES (?, DateTime('now'));

-- -- name: SearchSiteIndex :many
-- SELECT * FROM site_index 
--     WHERE text MATCH ?
--     ORDER BY rank
--     LIMIT ?;



-- -- name: SearchSiteIndexSnippetsWithPath :many
-- SELECT path, snippet(site_index, 2, '<b>', '</b>', '...', 16)
--     FROM site_index 
--     WHERE text MATCH ?
--     AND path LIKE ?
--     ORDER BY rank
--     LIMIT ?;

