-- migrate:up

-- there is no 'if not exists' for create type
-- https://stackoverflow.com/questions/7624919/check-if-a-user-defined-type-already-exists-in-postgresql
-- also, for future migrations:
-- https://stackoverflow.com/questions/1771543/adding-a-new-value-to-an-existing-enum-type

-------------------------------------------------
-- schemes
-------------------------------------------------
create table schemes (
  id integer not null,
  scheme text not null,
  unique(scheme)
);

insert into schemes 
	(id, scheme) 
	values 
	(1, 'https'), 
  (2, 'http'), 
  (3, 'ftp')
	on conflict do nothing
;

--------------------------------------------------
-- tlds
--------------------------------------------------
create table tlds (
  id integer not null,
  tld text not null,
  unique(id),
  unique(tld),
  constraint domain64_tld check (id > 0 and id < 256)
);


  -- gov: '01',
  -- mil: '02',
  -- com: '03',
  -- edu: '04',
  -- net: '05',
  -- org: '06',
insert into tlds
  (id, tld)
  values
  (1, 'gov'), 
  (2, 'mil'), 
  (3, 'com'), 
  (4, 'edu'),
  (5, 'net'),
  (6, 'org')
  on conflict do nothing
;

---------------------------------------------------
-- content type
---------------------------------------------------
create table content_types (
  id integer not null,
  content_type text not null
  unique(content_type)
);

insert into content_types
  (id, content_type)
  values
  (1, 'application/octet-stream'), 
  (2, 'text/html'), 
  (3, 'application/pdf')
  on conflict do nothing
;

create table tags (
  id integer not null,
  tag text not null,
  attr text null,
  unique(tag)
;

insert into tags
  (id, tag)
  values
  (1, 'title'),
  (2, 'p'),
  (3, 'div'),
  (4, 'h1'),
  (5,'h2'),
  (6,'h3'),
  (7,'h4'),
  (8,'h5'),
  (9,'h6'),
  (10,'a'),
  (11,'div'),
  (12,'th'),
  (13,'td')
  on conflict do nothing
;

insert into tags
  (id, tag, attr)
  values
  (14, 'img', 'alt')
  on conflict do nothing
;

create table hosts (
  id bigint generated always as identity primary key,
  next_fetch timestamp not null,
  domain64 bigint,
  host text,
  unique(id),
  unique(domain),
  constraint domain64_domain 
    check (domain64 > 0 and domain64 <= select(x'FFFFFFFFFFFFFFFF'))
);

create table guestbook (
  id bigint generated always as identity primary key,
  last_modified timestamp,
  last_fetched timestamp not null,
  next_fetch timestamp not null,
  scheme integer not null references schemes(id) default 1,
  host integer references hosts(id) not null,
  content_type integer references content_types(id) default 1,
  content_length integer not null default 0,
  path text not null,
  unique (host, path)
);



CREATE TABLE raw_content (
  id PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  -- domain64 TEXT NOT NULL,
  -- path TEXT NOT NULL,
  guestbook_id references guestbook(id) NOT NULL,
  tag TEXT default 'p',
  content TEXT NOT NULL
)

-- migrate:down
drop table if exists guestbook cascade;


