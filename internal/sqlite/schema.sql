CREATE TABLE metadata (
  id INTEGER PRIMARY KEY,
  column TEXT,
  value TEXT
);

INSERT INTO metadata 
  (column, value)
  VALUES
  ("version", "1"),
  ("datetime_created", DATETIME()),
  ("date_created", DATE())
;  

CREATE TABLE paths (
  id INTEGER PRIMARY KEY,
  path TEXT NOT NULL,
  UNIQUE(path)
);

CREATE TABLE titles (
  id INTEGER PRIMARY KEY,
  path_id INTEGER NOT NULL,
  kind INTEGER NOT NULL,
  txt TEXT NOT NULL,
  FOREIGN KEY(path_id) REFERENCES paths(id)
);

CREATE TABLE headers (
  id INTEGER PRIMARY KEY,
  path_id INTEGER NOT NULL,
  kind INTEGER NOT NULL,
  level INTEGER NOT NULL,
  txt TEXT NOT NULL,
  FOREIGN KEY(path_id) REFERENCES paths(id)
);

CREATE TABLE bodies (
  id INTEGER PRIMARY KEY,
  path_id INTEGER NOT NULL,
  kind INTEGER NOT NULL,
  tag TEXT NOT NULL,
  txt TEXT NOT NULL,
  FOREIGN KEY(path_id) REFERENCES paths(id)
);

CREATE VIRTUAL TABLE titles_fts USING fts5(
  path_id UNINDEXED,
  kind UNINDEXED,
  txt, 
  content='titles',
  content_rowid='id' 
);

CREATE VIRTUAL TABLE headers_fts USING fts5(
  path_id UNINDEXED,
  kind UNINDEXED,
  level UNINDEXED,
  txt,
  content='headers',
  content_rowid='id' 
);

CREATE VIRTUAL TABLE bodies_fts USING fts5(
  path_id UNINDEXED,
  kind UNINDEXED,
  txt,
  content='bodies',
  content_rowid='id' 
);

CREATE TRIGGER titles_ai AFTER INSERT ON titles
    BEGIN
        INSERT INTO titles_fts (path_id, kind, txt)
        VALUES (new.path_id, new.kind, new.txt);
    END;

CREATE TRIGGER headers_ai AFTER INSERT ON headers
    BEGIN
        INSERT INTO headers_fts (path_id, kind, txt, level)
        VALUES (new.path_id, new.txt, new.kind, new.level);
    END;

CREATE TRIGGER bodies_ai AFTER INSERT ON bodies
    BEGIN
        INSERT INTO bodies_fts (path_id, kind, txt)
        VALUES (new.path_id, new.kind, new.txt);
    END;
