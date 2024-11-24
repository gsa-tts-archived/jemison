
select count(*) from river_job where kind='fetch' and state='completed';

select count(*) from river_job where kind='extract';
select count(*) from river_job where kind='walk' and state='completed';

select (select count(*) from river_job where kind='fetch' and state='completed') /
(select extract(minute from (
	select maxf-minf from
		(select max(rj.finalized_at) as maxf, min(rj.finalized_at) as minf from river_job rj
		where kind='fetch')))) as fetch_jobs_per_minute,
	(select count(*) from river_job where kind='fetch' and state='completed') as fetch_total



select (select count(*) from river_job where kind='extract' and state='completed') /
(select extract(minute from (
	select maxf-minf from
		(select max(rj.finalized_at) as maxf, min(rj.finalized_at) as minf from river_job rj
		where kind='extract')))) as extract_jobs_per_minute,
	(select count(*) from river_job where kind='extract' and state='completed') as extract_total

select (select count(*) from river_job where kind='walk' and state='completed') /
(select extract(minute from (
	select maxf-minf from
		(select max(rj.finalized_at) as maxf, min(rj.finalized_at) as minf from river_job rj
		where kind='walk')))) as walk_jobs_per_minute,
		(select count(*) from river_job where kind='walk' and state='completed') as walk_total