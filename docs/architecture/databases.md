# databases

The engine is supported by multiple databases. We do this because different databases may want/need to scale differently in the future, and we are taking a preemptive step in that direction.

## unique data values

There are a small number of unique data representations in our application that merit note.

* the [domain64](domain64.md) representation for domain names. This is a 64-bit integer representation of domain names that we use for search optimization and table partitioning.

## queues

The queues database serves only one purpose: to handle the queues used by Jemison.

Our queueing system gets hit hard, and therefore we do all of that work on one database. Further, the table migrations are automanaged by the [river](https://riverqueue.com/) library. Keeping it separate both protects the queue tables as well as any operational tables we create in the application.

## work

The "work" database is where application tables specific to the processing of data live. 

### guestbook

The guestbook is where we keep track of URLs that have been/want to be searched. These tables live in the `cmd/migrate` app, which handles our migrations on every deploy. [These are dbmate migrations](https://github.com/GSA-TTS/jemison/tree/main/cmd/migrate/work_db/db/migrations).

```sql
create table guestbook (
  id bigint generated always as identity primary key,
  domain64 bigint not null,
  last_modified timestamp,
  last_fetched timestamp,
  next_fetch timestamp not null,
  scheme integer not null default 1,
  content_type integer not null default 1,
  content_length integer not null default 0,
  path text not null,
  unique (domain64, path)
);
```

The dates drive a significant part of the entree/fetch algorithms.

* `last_modified` is EITHER the timestamp provided by the remote webserver for any given page, OR if not present, we assign this value in `fetch`, setting it to the last fetched timestamp.
* `last_fetched` is the time that the page was fetched. This is updated every time we fetch the page.
* `next_fetch` is a computed value; if a page is intended to be fetched weekly, then `fetch` will set this as the current time plus one week at the time the page is fetched. 

### hosts

```sql
create table hosts (
  id bigint generated always as identity primary key,
  domain64 bigint,
  next_fetch timestamp not null,
  unique(id),
  unique(domain64),
  constraint domain64_domain 
    check (domain64 > 0 and domain64 <= max_bigint())
)
;
```

Like the `guestbook`, this table plays a role in determining whether a given domain should be crawled. If we want to crawl a domain *right now*, we set the `next_fetch` value in this table to yesterday, allowing all crawls of URLs under this domain to be valid. 

## search

The `search` database holds our data pipelines and the tables that get actively searched. 

This database is not (yet) well designed. Currently, there is a notion of a `raw_content` table, which is where `pack` deposits text.

```sql
CREATE TABLE raw_content (
  id BIGSERIAL PRIMARY KEY,
  host_path BIGINT references guestbook(id),
  tag TEXT default ,
  content TEXT 
)
```

From there, it is unclear how best to structure and optimize the content. 

There are two early-stage ideas. Both have tradeoffs in terms of performance and implementation complexity, and it is not clear yet which to pursue.


### one idea: inheritence.

https://www.postgresql.org/docs/current/tutorial-inheritance.html

We could define a searchable table as `gov`. 

```sql
create table gov (
  id ...,
  host_path ...,
  tag ...,
  content ...
);
```

From there, we could have *empty* inheritence tables.

```sql
create table gsa () inherits (gov);
create table hhs () inherits (gov);
create table nih () inherits (gov);
```

and, from there, the next level down:

```sql
create table cc () inherits (nih);
create table nccih () inherits (nih);
create table nia () inherits (nih);
```

Then, insertions happen at the **leaves**. That is, we only insert at the lowest level of the hierarchy. However, we can then query tables higher up, and get results from the entire tree.

This does two things:

1. It lets queries against a given domain happen naturally. If we want to query `nia.nih.gov`, we target that table with our query.
2. If we want to query all of `nih`, then we query the `nih` table.
3. If we want to query everything, we target `gov` (or another tld).

Given that we are going to treat these tables as build artifacts, we can always regenerate them. And, it is possible to add new tables through a migration easily; we just add a new create table statement.

(See [this article](https://medium.com/miro-engineering/sql-migrations-in-postgresql-part-1-bc38ec1cbe75) about partioning/inheritence, indexing, and migrations. It's gold.)

### declarative partitioning

Another approach is to use `PARTITION`s.

This would suggest our root table has columns we can use to drive the derivative partitions.

```sql
create table gov (
  id ...,
  domain64 BIGINT,
  host_path ...,
  tag ...,
  content ...
  partition by range(domain64)
);
```

To encode all of the TLDs, domains, and subdomains we will encounter, we'll use a `domain64` encoding. Why? It maps the entire URL space into a single, 64-bit number (or, `BIGINT`).

```
FF:FFFFFF:FFFFFF:FF
```

or

```
tld:domain:subdomain:subsub
```

This is described more in detail in [domain64.md](domain64.md).

As an example:

| tld | domain | sub |                  hex |               dec |
|-----|--------|-----|----------------------|-------------------|
| gov |    gsa |  _  |   #x0100000100000000 | 72057598332895232 |
| gov |    gsa | tts |   #x0100000100000100 | 72057598332895488 |
| gov |    gsa | api |   #x0100000100000200 | 72057598332895744 |

GSA is from the range #x0100000001000000 -> #x0100000001FFFFFF, or 72057594054705152 -> 72057594071482367 (a diff of 16777215). Nothing else can be in that range, because we're using the bitstring to partition off ranges of numbers.

Now, everything becomes bitwise operations on 64-bit integers, which will be fast everywhere... and, our semantics map well to our domain.

Partitioning to get a table with only GSA entries is

```sql
CREATE TABLE govgsa PARTITION OF gov
    FOR VALUES FROM (72057598332895232) TO (72057602627862527);
```

Or, just one subdomain in the space:

```sql
CREATE TABLE govgsatts PARTITION OF gov
    FOR VALUES FROM (72057598332895488) TO (72057598332895743);
```

or we can keep the hex representation:

```sql
CREATE TABLE govgsatts PARTITION OF gov
    FOR VALUES FROM (select x'0100000100000100') TO (select x'01000001000001FF');
```

All table operations are on the top-level table (insert, etc.), the indexes and whatnot are inherited automatically, and I can search the TLD, domain, or subdomain without difficulty---because it all becomes a question of what range the `domain64` value is in.


