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

create table searchable_content (
  domain64 bigint not null,
  path text not null,
  tag text default 'p' not null,
  content text not null,
  fts tsvector generated always as (to_tsvector('english', content)) STORED
)
;

-- alter table searchable_content
-- add column fts tsvector
-- generated always
-- as (to_tsvector('english', content)) stored;


create or replace function fun_raw_content_after_insert()
returns trigger as $$
begin
    insert into searchable_content 
      (domain64, path, tag, content)
    values 
      (new.domain64, new.path, new.tag, new.content);
    return new;
end;
$$ language plpgsql;

create trigger trigger_raw_content_after_insert
after insert on raw_content
for each row
execute function fun_raw_content_after_insert();

-- migrate:down

