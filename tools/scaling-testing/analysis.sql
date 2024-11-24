
select count(*) from river_job where kind='fetch' and state='completed';

select count(*) from river_job where kind='extract';
select count(*) from river_job where kind='walk' and state='completed';


-- entree JPM
select (select count(*) from river_job where kind='entree' and state='completed') /
(select extract(minute from (
	select maxf-minf from
		(select max(rj.finalized_at) as maxf, min(rj.finalized_at) as minf from river_job rj
		where kind='entree')))) as entree_jobs_per_minute,
	(select count(*) from river_job where kind='entree' and state='completed') as entree_total

-- fetch JPM
select (select count(*) from river_job where kind='fetch' and state='completed') /
(select extract(minute from (
	select maxf-minf from
		(select max(rj.finalized_at) as maxf, min(rj.finalized_at) as minf from river_job rj
		where kind='fetch')))) as fetch_jobs_per_minute,
	(select count(*) from river_job where kind='fetch' and state='completed') as fetch_total

-- extract JPM
select (select count(*) from river_job where kind='extract' and state='completed') 
	/
	(select extract(minute from 
		(select maxf-minf from
			(select max(rj.finalized_at) as maxf, min(rj.finalized_at) as minf from river_job rj
			where kind='extract')))) as extract_jobs_per_minute,
		(select count(*) from river_job where kind='extract' and state='completed') as extract_total

-- walk JPM
select (select count(*) from river_job where kind='walk' and state='completed') /
(select extract(minute from (
	select maxf-minf from
		(select max(rj.finalized_at) as maxf, min(rj.finalized_at) as minf from river_job rj
		where kind='walk')))) as walk_jobs_per_minute,
		(select count(*) from river_job where kind='walk' and state='completed') as walk_total

-- host overview
select 
	hname.host, gb.host, max(gb.last_fetched)-min(gb.last_fetched) as duration,
	count(gb.path) as count, 
	count(gb.path)/nullif(extract(minute from max(gb.last_fetched)-min(gb.last_fetched)), 0) as ppm
from guestbook gb, hosts hname
where hname.id = gb.host
group by gb.host, hname.host
order by count desc, ppm asc

-- overall ppm
select count(gb.path) / extract(minute from (select maxf-minf from
(select max(gb.last_fetched) as maxf, min(gb.last_fetched) as minf
from guestbook gb) as duration)) as ppm from guestbook gb

-- unique paths in guestbook
select count(gb.path) from guestbook gb

-- run duration based on guestbook
select extract(minute from (select maxf-minf from
(select max(gb.last_fetched) as maxf, min(gb.last_fetched) as minf
from guestbook gb) as duration)) minutes


-- Getting a summary of the job state
select kind, state, count(*) from river_job group by kind, state

-- PDFs in guestbook per domain
select h.host, count(*) as count, 'pdf' as type
	FROM hosts h, guestbook gb
	where 
		path ILIKE '%pdf'
		AND
		h.id = gb.host
	group by gb.host, h.host
	ORDER BY count DESC