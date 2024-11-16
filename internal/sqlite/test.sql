-- https://github.com/mattn/go-sqlite3?tab=readme-ov-file

CREATE TABLE paths (
  path TEXT,
  UNIQUE(path)
);

CREATE TABLE titles (
  path_rowid INTEGER,
  kind INTEGER,
  txt TEXT,
  FOREIGN KEY(path_rowid) REFERENCES paths(rowid)
);

CREATE TABLE headers (
  path_rowid INTEGER,
  kind INTEGER,
  level INTEGER,
  txt TEXT,
  FOREIGN KEY(path_rowid) REFERENCES paths(rowid)
);

CREATE TABLE bodies (
  path_rowid INTEGER,
  kind INTEGER,
  tag TEXT,
  txt TEXT,
  FOREIGN KEY(path_rowid) REFERENCES paths(rowid)
);

CREATE VIRTUAL TABLE titles_fts USING fts5(
  path_rowid UNINDEXED,
  kind UNINDEXED,
  txt, 
  content='titles',
  content_rowid='rowid' 
);

CREATE VIRTUAL TABLE headers_fts USING fts5(
  path_rowid UNINDEXED,
  kind UNINDEXED,
  level UNINDEXED,
  txt,
  content='headers',
  content_rowid='rowid' 
);

CREATE VIRTUAL TABLE bodies_fts USING fts5(
  path_rowid UNINDEXED,
  kind UNINDEXED,
  txt,
  content='bodies',
  content_rowid='rowid' 
);

CREATE TRIGGER titles_ai AFTER INSERT ON titles
    BEGIN
        INSERT INTO titles_fts (rowid, path_rowid, kind, txt)
        VALUES (new.rowid, new.path_rowid, new.kind, new.txt);
    END;

CREATE TRIGGER headers_ai AFTER INSERT ON headers
    BEGIN
        INSERT INTO headers_fts (rowid, path_rowid, kind, txt, level)
        VALUES (new.rowid, new.path_rowid, new.txt, new.kind, new.level);
    END;

CREATE TRIGGER bodies_ai AFTER INSERT ON bodies
    BEGIN
        INSERT INTO bodies_fts (rowid, path_rowid, kind, txt)
        VALUES (new.rowid, new.path_rowid, new.kind, new.txt);
    END;

insert into paths 
  (path)
  values ("/a"), ("/b");

begin;
insert into titles 
  (path_rowid, kind, txt)
  values 
    (1, 0, "T dogs of the world"),
    (2, 0, "T an incremental approach")
;
commit;

begin;
insert into headers
  (path_rowid, kind, level, txt)
  values
    (1, 1, 1, "H the outline"),
    (1, 1, 1, "H a shaggy dog story"),
    (1, 1, 2, "H the cats and dogs"),
    (2, 1, 1, "H something new"),
    (2, 1, 2, "H something old"),
    (2, 1, 3, "H something borrowed")
;
commit;

begin;
insert into bodies
  (path_rowid, kind, tag, txt)
  values
    (1, 2, "p", "C the quick brown fox jumped over the lazy dog"),
    (1, 2, "p", "C something in the way she moves"),
    (1, 2, "li", "C get milk when you go to the grocery"),
    (2, 2, "p", "C a typewriter is a clacky thing"),
    (2, 2, "p", "C don't stand so close to me"),
    (2, 2, "li", "C get the butter when you go to the grocery")
;
commit;

SELECT PRINTF("dog");
SELECT PRINTF("---");
select * from bodies_fts where bodies_fts match 'the';
SELECT PRINTF("");

SELECT PRINTF("headers");
SELECT PRINTF("---");
select * from headers_fts where headers_fts match 'the';
SELECT PRINTF("");

SELECT PRINTF("contents");
SELECT PRINTF("---");
select * from bodies_fts where bodies_fts match 'the'
order by rank asc;
SELECT PRINTF("");

-- An integer indicating the index of the FTS table column to select the returned text from. Columns are numbered from left to right starting at zero. A negative value indicates that the column should be automatically selected.
-- The text to insert before each phrase match within the returned text.
-- The text to insert after each phrase match within the returned text.
-- The text to add to the start or end of the selected text to indicate that the returned text does not occur at the start or end of its column, respectively.
-- The maximum number of tokens in the returned text. This must be greater than zero and equal to or less than 64.

SELECT PRINTF("just content");
SELECT PRINTF("---");
SELECT 1 as weight, rowid, txt, path_rowid, rank
  FROM bodies_fts
  WHERE bodies_fts MATCH 'the'
ORDER BY weight DESC, rank ASC
;
SELECT PRINTF("");

SELECT PRINTF("raw query");
SELECT PRINTF("---");
SELECT path_rowid, 4.0 as weight, rank,  txt
  FROM titles_fts
  WHERE titles_fts MATCH 'the'
UNION ALL
SELECT path_rowid, 2.0 as weight, rank, txt
  FROM headers_fts
  WHERE headers_fts MATCH 'the'
UNION ALL
SELECT path_rowid, 1.0 as weight, rank, txt
  FROM bodies_fts
  WHERE bodies_fts MATCH 'the'
ORDER BY weight DESC, rank ASC
;
SELECT PRINTF("");

SELECT PRINTF("weighted query");
SELECT PRINTF("---");

SELECT 
  path_rowid,
  (SELECT path from paths WHERE rowid = path_rowid), 
  kind,
  weight,
  rank,
  txt
  FROM
    (SELECT path_rowid, 4.0 as weight, rank, kind, txt
      FROM titles_fts
      WHERE titles_fts MATCH 'the'
    UNION ALL
    SELECT path_rowid, 2.0 as weight, rank, kind, txt
      FROM headers_fts
      WHERE headers_fts MATCH 'the'
    UNION ALL
    SELECT path_rowid, 1.0 as weight, rank, kind, txt
      FROM bodies_fts
      WHERE bodies_fts MATCH 'the'
    ORDER BY weight DESC, rank ASC)
  as the_query
  --GROUP BY path_rowid
  ORDER BY rank DESC
  ;
