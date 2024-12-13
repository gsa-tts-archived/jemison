-- migrate:up

-- there is no 'if not exists' for create type
-- https://stackoverflow.com/questions/7624919/check-if-a-user-defined-type-already-exists-in-postgresql
-- also, for future migrations:
-- https://stackoverflow.com/questions/1771543/adding-a-new-value-to-an-existing-enum-type


create or replace function public.max_bigint()
  returns bigint
  language sql immutable parallel safe as
--'select 18446744073709551615'; -- (sub1 (expt 2 64))
-- These are signed integers...
'select 9223372036854775807';

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



CREATE TABLE raw_content (
  id bigint generated always as identity primary key,
  guestbook_id bigint references guestbook(id) NOT NULL,
  tag integer default 1,
  content TEXT NOT NULL
)

-- migrate:down
drop table if exists guestbook cascade;


