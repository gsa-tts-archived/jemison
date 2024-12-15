-- migrate:up
create table if not exists metadata (
  id integer generated always as identity primary key,
  col text,
  value text
)
;

insert into metadata 
  (col, value)
  values
  ('version', '1'),
  ('datetime_created', now()::text),
  ('date_created', now()::date)
;  

create table raw_content (
  id bigint generated always as identity primary key,
  domain64 bigint not null,
  path text not null,
  tag text default 'p' not null,
  content text not null
)
;

-- migrate:down

