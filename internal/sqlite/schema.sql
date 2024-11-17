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
  title TEXT NOT NULL,
  FOREIGN KEY(path_id) REFERENCES paths(id)
);

CREATE TABLE headers (
  id INTEGER PRIMARY KEY,
  path_id INTEGER NOT NULL,
  kind INTEGER NOT NULL,
  level INTEGER NOT NULL,
  header TEXT NOT NULL,
  FOREIGN KEY(path_id) REFERENCES paths(id)
);

CREATE TABLE bodies (
  id INTEGER PRIMARY KEY,
  path_id INTEGER NOT NULL,
  kind INTEGER NOT NULL,
  tag TEXT NOT NULL,
  body TEXT NOT NULL,
  FOREIGN KEY(path_id) REFERENCES paths(id)
);

CREATE VIRTUAL TABLE titles_fts USING fts5(
  path_id UNINDEXED,
  kind UNINDEXED,
  title,
  content='titles',
  content_rowid='id' 
);

CREATE VIRTUAL TABLE headers_fts USING fts5(
  path_id UNINDEXED,
  kind UNINDEXED,
  level UNINDEXED,
  header,
  content='headers',
  content_rowid='id' 
);

CREATE VIRTUAL TABLE bodies_fts USING fts5(
  path_id UNINDEXED,
  kind UNINDEXED,
  body,
  content='bodies',
  content_rowid='id' 
);

CREATE TRIGGER titles_ai AFTER INSERT ON titles
    BEGIN
        INSERT INTO titles_fts (path_id, kind, title)
        VALUES (new.path_id, new.kind, new.title);
    END;

CREATE TRIGGER headers_ai AFTER INSERT ON headers
    BEGIN
        INSERT INTO headers_fts (path_id, kind, level, header)
        VALUES (new.path_id, new.kind, new.level, new.header);
    END;

CREATE TRIGGER bodies_ai AFTER INSERT ON bodies
    BEGIN
        INSERT INTO bodies_fts (path_id, kind, body)
        VALUES (new.path_id, new.kind, new.body);
    END;
