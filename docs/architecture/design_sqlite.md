# sqlite / `serve`

The output of `pack` (and `superpack`) is an SQLite file.

We copy this file from S3 to a host, and serve queries out of it directly. This is because it is *very* fast, and it saves us from having to garden/care for a live database server under production loads. (Or: our search can be scaled horizontally across small EC2 instances, as opposed to trying to figure out how to get read replicas/etc. running well under the brokered environment of cloud.gov).

## connection string / driver

https://github.com/mattn/go-sqlite3?tab=readme-ov-file#:~:text=see%20PRAGMA%20defer_foreign_keys-,Foreign%20Keys,-_foreign_keys%20%7C%20_fk

We need foreign keys.

`_fk=true`


## metadata

Every SQLite file needs a `metadata` table. This helps us know what we are serving from. For example, we may someday go from "version 1" of our tables to "version 2." When we do this:

* `pack` will start producing version 2 tables, and
* `serve` will need to know how to serve both v1 and v2 tables.

This will only be true until we have re-crawled or otherwise regenerated all of our SQLite databases. At that point, `serve` can "forget" how to handle v1 tables.

### `metadata`

Metadata is a k/v table. Therefore, the design is two columns. Remember, SQLite *barely* has types.

| column | type |
| --- | --- |
| key | TEXT |
| value | TEXT |

The following keys MUST exist:
* `version`, which is an integer. Only update for breaking changes.
* `last_updated`, a `DATE`. This is the date we wrote the file.

### `paths`

These are not for searching, but to keep tables small.

```sql
CREATE TABLE paths (
  id PRIMARY KEY,
  path TEXT,
  UNIQUE(path)
);
```

### `titles`

Titles will go in their own table.

```sql
CREATE TABLE titles (
  id PRIMARY KEY,
  path INTEGER,
  title TEXT,
  FOREIGN KEY(path) REFERENCES paths(id)
);
```

### `headers`

The `level` is the integer portion of `H1`, `H2`, etc.

```sql
CREATE TABLE headers (
  id PRIMARY KEY,
  path INTEGER,
  level INTEGER,
  header TEXT,
  FOREIGN KEY(path) REFERENCES paths(id)
);
```

### `content`

This is the body content of the page. So, everything but the other tables.

```sql
CREATE TABLE contents (
  id PRIMARY KEY,
  path INTEGER,
  tag TEXT,
  content TEXT,
  FOREIGN KEY(path) REFERENCES paths(id)
);
```

## FTS5 tables

```sql
CREATE VIRTUAL TABLE titles_fts USING fts5(
    title, 
    content='titles',
    content_rowid='id' 
);

CREATE VIRTUAL TABLE headers_fts USING fts5(
    header,
    content='header',
    content_rowid='id' 
);

CREATE VIRTUAL TABLE contents_fts USING fts5(
    content,
    content='contents',
    content_rowid='id' 
);
```


### FTS5 triggers

We need insert triggers for this model.

```sql
CREATE TRIGGER titles_ai AFTER INSERT ON titles
    BEGIN
        INSERT INTO titles_fts (rowid, title)
        VALUES (new.id, new.title);
    END;

CREATE TRIGGER headers_ai AFTER INSERT ON headers
    BEGIN
        INSERT INTO headers_fts (rowid, header)
        VALUES (new.id, new.header);
    END;

CREATE TRIGGER contents_ai AFTER INSERT ON contents
    BEGIN
        INSERT INTO contents_fts (rowid, content)
        VALUES (new.id, new.content);
    END;
```