# jemison


`jemison` is a prototype for a small, scalable, maintainable search platform. 

It is based on a series of experiments in search and embedded database technologies.

* Experiment four explored the performance of SQLite databases for full-text retrieval via HTTP as well as comparisons between compressed and uncompressed databases.
* Experiment [six](https://github.com/jadudm/six) explored the creation of an end-to-end web crawler/indexer/search server in 2000 lines of Go.
* Experiment [eight](https://github.com/GSA-TTS/jemison) revisited this idea, focusing on the idea of search as an *observable data pipeline*. 

`jemison` is named for the pioneering and innovating explorer of medicine and space, Dr. Mae Jemison, the first African-American woman in space. 

## running the experiment

In the top directory, first build the base/build container:

```
make docker
```

Then, run the stack.

```
make up
```

The `up` compiles all of the services, generates the database API, unpacks the USWDS assets into place, and launches the stack. This kind of build is possible when working on an Intel-based platform.

## containerized build

```
make docker
```

to start, and 

```
make macup
```

to build and run the applications within a containerized environment. This runs the build in a Linux context (against the image used by CF/cloud.gov). This style build is required when running on non-Intel platforms.

## interacting with components

An API key is hard-coded into the `compose.yaml` file; it is `lego`. Obviously, this is not suitable for production use.

To begin a crawl:

```
http put http://localhost:10001/api/entree/full/pass \
  scheme=https \
  path="/" api-key=lego \
  host=digitalcorps.gsa.gov
```

The URL parameters to the `admin` component determine if we are doing a full crawl and whether or not we have a hall pass. A "hall pass" lets us fetch a page even if the deadline for the next fetch has not passed. This is for emergency updating of pages/domains outside of a normal scheudle.

* `/single/no` fetches a single page, with no pass; if the deadline has not been met, we don't fetch the page.
* `/single/pass` fetches the page immediately
* `/full/no` does a full crawl, but no hall pass is issued. This will check the deadlines on all existing pages, and we will only queue fetches on pages whose deadlines have been passed.
* `/full/pass` will set the crawl deadlines on all known pages to *yesterday*, and then kick off a crawl. 

When resetting the schedule (`/full/pass`), we set a new deadline in keeping with the site's normal schedule. For example, if a site is scheduled to run weekly, then we will, as we fetch, set the deadline on a new fetch to one week in the future. In this regard, a weekly schedule is not fixed---it is not necessarily "every Monday." Instead, it is one week from the last time it was fetched. (This may turn out to be confusing behavior...)

### configuration

All configuration happens through JSonnet files in `config`. The YAML/JSON files are generated from these sources.

* Files in the `services` directory control the golang services (`entree`, `fetch`, etc.). 
* Files in the `domain64` directory control the domains we work with.
* The `allowed_hosts` file limits what can be crawled. 

The `domain64` directory should be edited with caution; renumbering, once in production, would cause no end of confusion, and end up returning results from the wrong sites. These files are slow to generate, but we're willing to pay the price for the correctness that JSonnet helps guarantee.

## searching 

After a site is walked and packed, the data is stored in Postgres for searching. The `serve` component transforms the API into search queries, and also provides a rudimentary template simulating a search result landing page (SERP).

```
http post http://localhost:10000/api/search \
  host=www.fac.gov \
  d64_start="0" \
  d64_end="143833713099145216" \
  terms="audit compliance" 
```

(At the time of writing, the `host` parameter may be redundant...)

The domain64 start and end values control what part of the domain space to search. See the docs on `domain64` for more information.

A [WWW-based search interface](http://localhost:10000/gov) can be found at [http://localhost:10000/gov](http://localhost:10000/gov). Note that the URL

```
http://localhost:10004/search/{TLD}/{DOMAIN}/{SUBDOMAIN}
```

determines what indexed content will be searched. (E.g. if you have indexed `alice.gov` and `bob.gov`, selecting that database will navigate you to a URL like `http://localhost:10004/search/gov/alice`.)

This lets you search

* http://localhost:10000/gov
* http://localhost:10000/gov/nasa
* http://localhost:10000/gov/nasa/blogs

if you have crawled blogs.nasa.gov. A search for the TLD only will run a search against all content in that TLD.

## browsing the backend

The goal is to minimize required services. This stack *only* uses Postgres and S3. 

## browsing S3

The S3 filestore is simulated when running locally using a containerized version of [Minio](https://min.io).

![alt text](docs/images/minio.png)

Point a browser at [localhost:9001](http://localhost:9001) with the credentials `numbernine/numbernine` to browse.

### watching the queue

There is a UI for monitoring the queues.

![alt text](docs/images/riverui.png)

This lets you watch the queues at [localhost:11111](http://localhost:11111) provided by [River](https://riverqueue.com/), a queueing library/system built on Postgres. 


[pgweb](https://sosedoff.github.io/pgweb/) is included in the container stack for browsing the queue database directly and, if needed, editing. Pointign a browser at [localhost:22222](http://localhost:22222) will bring up `pgweb`. 

If you are running, and want to simulate total queue loss, run

```
truncate table river_job; truncate table river_leader; truncate table river_queue;
```

This will not break the app; it will, however, leave all of the services with nothing to do.

Alternatively, run

```
tools/wipe-queue.bash
```

which will do the same thing

## viewing the second database

[localhost:22223](http://localhost:22223) is a second `pgweb` instance that looks at the "working" database (`jemison-work-db`). This contains the guestbook (which tracks what URLs we've visisted) and other data that we need for running a crawler for the `.gov` domain. 

In generall, the DB should be used sparingly. In the case of the crawl loop (e.g. `entree -> fetch -> walk -> entree...`), we need to keep track of roughly 25M URLs (perhaps more?). And, because cloud.gov limits us to a 7GB disk, we cannot do this locally in SQLite---a database of 25M URLs is probably going to push (or exceed) our disk space.

So, we use Postgres for things that we know might grow beyond 6GB of storage. We trade some performance, but gain space. 

## other utilities

To run the stack without the services (just the backend of `minio`, `postgres`, and the UIs)

```
docker compose -f backend.yaml up
```

or 

```
make backend
```

## in case of emergency

What happens if there is a fire in production?

From a service/serving the public perspective, our databases are in S3. When the `serve` component(s) wake up, they copy their databases out of S3, and serve traffic. In this regard, we can probably always restart the production environment and serve the current state of our crawling efforts.

We are somewhat safe from a fire with regards to the crawler. That is, if we had to completely wipe both `queue-db` and `work-db`, we would be starting from a clean slate. The crawl sequence of `entree->fetch->walk->...` would have no prior state. This would be equivalent to believing that we had never crawled anything, ever.

One solution/possibility is that we could have a process that 

1. Looks at everything in the `fetch` bucket, and 
2. Inserts it into the `guestbook`, with sensible timestamps.

We could, by the same measure, load each SQLite database from the `serve` bucket, and rebuild our `guestbook` the same way. The `last_fetched` date can be pulled from the `fetch` metadata... meaning the `guestbook` can be rebuilt for the cost of ~25M S3 queries.

We can probably achieve something similar by pulling each `serve` database from the S3 bucket one at a time, and using the crawled URLs from the search DBs to rebuild the `guestbook`.

Either way, it is possible to rebuild the state of the system from the artifacts in S3.

If we rebuild *nothing*, `entree`* would enqueue all of the hosts in the dark of the night. Because we would have no cached URLs, it would be equivalent to kicking off all of our crawls at the same time. (Or, enqueueing them all at once.) At that point, we'd simply be re-crawling all of our known sites. It would take time, but we would (essentially) not be at a loss.

In the case that we lose *all* of our services (e.g. S3 is wiped), we would have to recrawl, and service would be interrupted until we could re-build the SQLite databases. This would probably take a month.

### conclusion

We should back up the SQLite databases periodically, as they are our disaster recovery path. Being able to restore the 3000 or so SQLite databases is what lets us serve results.

## by the numbers

The original "experiment number eight" came in at 2500 lines of code.

The expansion of services (e.g. breaking out hit tracking into `entree`, the addition of `validate`, etc.) has pushed the system from 2500 lines of Go to just over 4000. 

*The system was always going to grow as we head to production.* However, it is also still much, much smaller than the previous system, and remains architecturally cleaner. :shrug:

```
cloc --exclude-ext=yml,yaml,html,css,less,js,json,svg,scss,csv,jsonnet,libsonnet --fullpath --not-match-d=terraform/\.terraform --not-match-d=config --not-match-d=venv .
```

```
--------------------------------------------------------------------------------
Language                      files          blank        comment           code
--------------------------------------------------------------------------------
Go                               72            905            584           4295
JSON                              2              0              0           2852
Markdown                         28            492              0            981
make                             18             93              7            366
HCL                               5             49             37            300
SQL                               6             83            107            208
Python                            3             27             27            124
Dockerfile                        6             30             27             67
Bourne Shell                      6             11              0             33
Bourne Again Shell                1              2              0              6
NAnt script                       1              2              0              3
--------------------------------------------------------------------------------
SUM:                            148           1694            789           9235
--------------------------------------------------------------------------------
```

### local load testing


With a connection pool limit of 100, and no particular optimizations on the API (running in `debug` mode), we see roughly 200 req/s. This comes out to around 12K req/minute, which is twice what Search.gov sees around peak (5K req/minute).

These numbers are against a local stack, and therefore have no network overhead.

More testing is needed to see if there are failure modes here. It seems, from the test, that we returned 200 to everything, but some requests pushed out to 4s. At 2x peak, a 4s query response is *ok*. We can also do more to load balance, have read replicas, etc. if we start seeing  traffic scale to this level. We might also consider something like [pgbouncer](https://www.pgbouncer.org/).


```
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: search-stressor.js
        output: -

     scenarios: (100.00%) 4 scenarios, 1249 max VUs, 3m10s max duration (incl. graceful stop):
              * one_is_such_a_lonely_number: 1.00 iterations/s for 10s (maxVUs: 10-100, gracefulStop: 30s)
              * the_marching_100: 100.00 iterations/s for 30s (maxVUs: 100-200, startTime: 10s, gracefulStop: 30s)
              * a_new_millenium: 500.00 iterations/s for 30s (maxVUs: 100-1000, startTime: 40s, gracefulStop: 30s)
              * teeter_totter: Up to 750 looping VUs for 1m30s over 4 stages (gracefulRampDown: 10s, startTime: 1m10s, gracefulStop: 30s)


     ✓ status was 200

     checks.........................: 100.00% 36643 out of 36643
     data_received..................: 40 MB   246 kB/s
     data_sent......................: 7.4 MB  45 kB/s
     dropped_iterations.............: 1705    10.429829/s
     http_req_blocked...............: avg=25.59µs  min=1.15µs med=4.67µs  max=24.11ms p(90)=6.49µs  p(95)=9.74µs  
     http_req_connecting............: avg=17.73µs  min=0s     med=0s      max=24.06ms p(90)=0s      p(95)=0s      
   ✗ http_req_duration..............: avg=611.99ms min=1.14ms med=91.17ms max=37.44s  p(90)=1.22s   p(95)=3.56s   
       { expected_response:true }...: avg=611.99ms min=1.14ms med=91.17ms max=37.44s  p(90)=1.22s   p(95)=3.56s   
   ✓ http_req_failed................: 0.00%   0 out of 36643
     http_req_receiving.............: avg=68.12µs  min=7.37µs med=37.35µs max=22.92ms p(90)=74.34µs p(95)=158.44µs
     http_req_sending...............: avg=24.27µs  min=3.38µs med=16.18µs max=22.84ms p(90)=23.65µs p(95)=31.11µs 
     http_req_tls_handshaking.......: avg=0s       min=0s     med=0s      max=0s      p(90)=0s      p(95)=0s      
     http_req_waiting...............: avg=611.89ms min=1.08ms med=91.12ms max=37.44s  p(90)=1.22s   p(95)=3.56s   
     http_reqs......................: 36643   224.152622/s
     iteration_duration.............: avg=1.61s    min=1s     med=1.09s   max=38.44s  p(90)=2.22s   p(95)=4.56s   
     iterations.....................: 36642   224.146505/s
     vus............................: 14      min=1              max=933 
     vus_max........................: 1249    min=750            max=1249
```

Points of comparison:

* https://stackoverflow.com/questions/373098/whats-the-average-requests-per-second-for-a-production-web-application
* https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/
* https://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/