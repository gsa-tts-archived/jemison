-- migrate:up

-- there is no 'if not exists' for create type
-- https://stackoverflow.com/questions/7624919/check-if-a-user-defined-type-already-exists-in-postgresql
-- also, for future migrations:
-- https://stackoverflow.com/questions/1771543/adding-a-new-value-to-an-existing-enum-type


create function public.max_bigint()
  returns bigint
  language sql immutable parallel safe as
'select 18446744073709551615';

create table hosts (
  id bigint generated always as identity primary key,
  next_fetch timestamp not null,
  domain64 bigint,
  host text,
  unique(id),
  unique(domain),
  constraint domain64_domain 
    check (domain64 > 0 and domain64 <= max_bigint())
)
;

create table guestbook (
  id bigint generated always as identity primary key,
  last_modified timestamp,
  last_fetched timestamp not null,
  next_fetch timestamp not null,
  host bigint references hosts(id) not null,
  scheme integer not null default 1,
  content_type integer default 1,
  content_length integer not null default 0,
  path text not null,
  unique (host, path)
);



CREATE TABLE raw_content (
  id bigint generated always as identity primary key,
  -- domain64 TEXT NOT NULL,
  -- path TEXT NOT NULL,
  guestbook_id bigint references guestbook(id) NOT NULL,
  tag integer default 1,
  content TEXT NOT NULL
)

-- migrate:down
drop table if exists guestbook cascade;


