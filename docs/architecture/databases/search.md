
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


