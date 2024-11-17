# jemison


`jemison` is a prototype for a small, scalable, maintainable search platform. 

It is based on a series of experiments in search and embedded database technologies.

* Experiment four explored the performance of SQLite databases for full-text retrieval via HTTP as well as comparisons between compressed and uncompressed databases.
* Experiment [six](https://github.com/jadudm/six) explored the creation of an end-to-end web crawler/indexer/search server in 2000 lines of Go.
* Experiment [eight](https://github.com/GSA-TTS/jemison) revisited this idea, focusing on the idea of search as an *observable data pipeline*. 

`jemison` is named for the pioneering and innovating explorer of medicine and space, Dr. Mae Jemison, the first African-American woman in space. 


> [!IMPORTANT]
> 
> `jemison` is:
> 
> 1. **MAINTAINABLE**: 15x fewer lines of code
> 2. **AFFORDABLE** 10x less resource intenstive (read: hosting costs)
> 3. **FAST**: 10x more performant (and scales horitzontally)

Why `jemison` is an MVP candidate for search:

1. `jemison` can already crawl and index HTML and PDFs, handling diacritics (e.g. Español).
2. `jemison` can deploy a full crawl/index/serve solution in 1GB of RAM on [cloud.gov](https://cloud.gov/)<sup>1</sup>.
3. `jemison` is 4000 lines of code: easy to maintain, with an extensible, distributed architecture.
4. `jemison` can handle 55K requests per minute (10x the current search.gov load) with 256**M**B of RAM.

<sup>1</sup> <em>This was true at one point. It needs a touch more now.</em>

## running the experiment

In the top directory, first build the base/build container:

```
make docker
```

Then, run the stack.

```
make up
```

The `up` compiles all of the services, generates the database API, unpacks the USWDS assets into place, and launches the stack.

## containerized build

```
make docker
```

to start, and 

```
make macup
```

to build and run the applications within a containerized environment. This runs the build in a Linux context (against the image used by CF/cloud.gov).

## interacting with components

An API key is hard-coded into the `compose.yaml` file; it is `lego`. Obviously, this is not suitable for production use.

To begin a crawl:

```
http put http://localhost:10001/api/entree/full/pass \
  scheme=https \
  host=digitalcorps.gsa.gov \
  path="/" api-key=lego
```

The URL parameters to the `admin` component determine if we are doing a full crawl and whether or not we have a hall pass. A "hall pass" lets us fetch a page even if the deadline for the next fetch has not passed. This is for emergency updating of pages/domains outside of a normal scheudle.

* `/single/no` fetches a single page, with no pass; if the deadline has not been met, we don't fetch the page.
* `/single/pass` fetches the page immediately
* `/full/no` does a full crawl, but no hall pass is issued. This will check the deadlines on all existing pages, and we will only queue fetches on pages whose deadlines have been passed.
* `/full/pass` will set the crawl deadlines on all known pages to *yesterday*, and then kick off a crawl. 

When resetting the schedule (`/full/pass`), we set a new deadline in keeping with the site's normal schedule. For example, if a site is scheduled to run weekly, then we will, as we fetch, set the deadline on a new fetch to one week in the future. In this regard, a weekly schedule is not fixed---it is not necessarily "every Monday." Instead, it is one week from the last time it was fetched. (This may turn out to be confusing behavior...)

The file `container.yaml` is a configuration file with a *few* end-user tunables. The FAC website is small in HTML, but is *large* because it contains 4 PDFs at ~2000 pages each. If you only want to index the HTML, set `extract_pdf` to `false`. (This is good for demonstration purposes.)

To fetch a single PDF and see it extracted:

```
http put http://localhost:10001/api/entree/full/pass \
  scheme=https \
  host=app.fac.gov \
  path="/dissemination/report/pdf/2023-09-GSAFAC-0000063050" \
  api-key=lego
```

(approximately 100 pages)

Note the `host` must be in the known hosts config (`schedule.jsonnet`).

## searching 

After a site is walked and packed, an SQLite file with full-text capabilities is generated. The `serve` component watches for completed files, grabs them from S3, and serves queries from the resulting SQLite database.

```
http POST localhost:10004/serve  host=fac.gov terms="community grant"
```

is how to search using the API; search terms are a single list, and SQLite pulls them apart.

A [WWW-based search interface](http://localhost:10004/search/) can be found at [http://localhost:10004/search/](http://localhost:10004/search/). You will be redirected to a page that lists all of the sites currently indexed. Note that the final part of the URL

```
http://localhost:10004/search/{HOST}
```

determines what indexed database will be searched. (E.g. if you have indexed `alice.gov` and `bob.gov`, selecting that database will navigate you to a URL like `http://localhost:10004/search/alice.gov`.)

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

To 

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
cloc --exclude-ext=yml,yaml,html,css,less,js,json,svg,scss --fullpath --not-match-d=terraform/\.terraform .
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

Using k6, I ran queries against random databases with random queries of length 1-4. This is with the script running locally against localhost, where network effects are minimal. Ramping up to 1000 simultaneous users, we see a throughput of 930 requests per second. 

That means that the search server can handle 55K (fifty-five thousand) requests per minute. The server in question **only has 512M of memory**. Scaling is linear with RAM, meaning that to double the number of simultaneous connections, we double the RAM.

search.gov sees approximately 5K requests per minute; with a tiny fraction of the resources of the current infrastructure, this design can handle 10x the load without breaking a sweat.


```
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: search-stressor.js
        output: -

     scenarios: (100.00%) 4 scenarios, 1999 max VUs, 4m10s max duration (incl. graceful stop):
              * one_is_such_a_lonely_number: 1.00 iterations/s for 10s (maxVUs: 10-100, gracefulStop: 30s)
              * the_marching_100: 100.00 iterations/s for 30s (maxVUs: 100-200, startTime: 10s, gracefulStop: 30s)
              * a_new_millenium: 1000.00 iterations/s for 1m0s (maxVUs: 100-1500, startTime: 40s, gracefulStop: 30s)
              * teeter_totter: Up to 1200 looping VUs for 2m0s over 4 stages (gracefulRampDown: 10s, startTime: 1m40s, gracefulStop: 30s)

     ✗ status was 200
      ↳  99% — ✓ 135719 / ✗ 395

     checks.........................: 99.70% 135719 out of 136114
     data_received..................: 118 MB 535 kB/s
     data_sent......................: 21 MB  97 kB/s
     dropped_iterations.............: 1170   5.294314/s
     http_req_blocked...............: avg=12.52µs min=0s       med=3.38µs  max=15.48ms p(90)=10.73µs p(95)=13.61µs 
     http_req_connecting............: avg=5.7µs   min=0s       med=0s      max=15.41ms p(90)=0s      p(95)=0s      
   ✓ http_req_duration..............: avg=27.51ms min=0s       med=8.1ms   max=1.82s   p(90)=64.57ms p(95)=94.26ms 
       { expected_response:true }...: avg=27.41ms min=273.16µs med=8.13ms  max=1.82s   p(90)=64.46ms p(95)=93.81ms 
   ✓ http_req_failed................: 0.29%  395 out of 136114
     http_req_receiving.............: avg=41.6µs  min=0s       med=25.51µs max=11.79ms p(90)=82.76µs p(95)=118.74µs
     http_req_sending...............: avg=18.56µs min=0s       med=10.18µs max=15.44ms p(90)=37.8µs  p(95)=48.56µs 
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s      p(90)=0s      p(95)=0s      
     http_req_waiting...............: avg=27.45ms min=0s       med=8.03ms  max=1.82s   p(90)=64.48ms p(95)=94.18ms 
     http_reqs......................: 136114 615.923323/s
     iteration_duration.............: avg=1.02s   min=1s       med=1s      max=2.83s   p(90)=1.06s   p(95)=1.09s   
     iterations.....................: 136114 615.923323/s
     vus............................: 4      min=1                max=1199
     vus_max........................: 1999   min=1200             max=1999


running (3m41.0s), 0000/1999 VUs, 136114 complete and 0 interrupted iterations
one_is_such_a_lonely_number ✓ [======================================] 000/010 VUs    10s   1.00 iters/s
the_marching_100            ✓ [======================================] 000/104 VUs    30s   100.00 iters/s
a_new_millenium             ✓ [======================================] 0000/1054 VUs  1m0s  1000.00 iters/s
teeter_totter               ✓ [======================================] 0000/1200 VUs  2m0s 
```

Points of comparison:

* https://stackoverflow.com/questions/373098/whats-the-average-requests-per-second-for-a-production-web-application
* https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/
* https://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/