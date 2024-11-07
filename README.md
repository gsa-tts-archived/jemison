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

The `up` compiles all of the services, generates the database API, unpacks the USWDS assets into place, and launches the stack.

## interacting with components

An API key is hard-coded into the `compose.yaml` file; it is `lego`. Obviously, this is not suitable for production use.

To begin a crawl:

```
http PUT localhost:10000/fetch host=fac.gov path=/ api_key=lego
```

The file `container.yaml` is a configuration file with a *few* end-user tunables. The FAC website is small in HTML, but is *large* because it contains 4 PDFs at ~2000 pages each. If you only want to index the HTML, set `extract_pdf` to `false`. (This is good for demonstration purposes.)

To fetch a single PDF and see it extracted:

```
http PUT localhost:10000/fetch host=app.fac.gov path=/dissemination/report/pdf/2023-09-GSAFAC-0000063050
```

(approximately 100 pages)

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

Point a browser at [localhost:9001](http://localhost:9001) with the credentials `nutnutnut/nutnutnut` to browse.

### watching the queue

There is a UI for monitoring the queues.

![alt text](docs/images/riverui.png)

This lets you watch the queues at [localhost:11111](http://localhost:11111) provided by [River](https://riverqueue.com/), a queueing library/system built on Postgres. 

### observing the database

[pgweb](https://sosedoff.github.io/pgweb/) is included in the container stack for browsing the database directly and, if needed, editing. Pointign a browser at [localhost:22222](http://localhost:22222) will bring up `pgweb`.

If you are running, and want to simulate total queue loss, run

```
truncate table river_job; truncate table river_leader; truncate table river_queue;
```

This will not break the app; it will, however, leave all of the services with nothing to do.

## other utilities

To run the stack without the services (just the backend of `minio`, `postgres`, and the UIs)

```
docker compose -f backend.yaml up
```

To 

## by the numbers

```
 docker run --rm -v ${PWD}:/tmp aldanial/cloc --exclude-dir=assets --fullpath --not-match-d=terraform/zips/* --not-match-d=terraform/app/* --not-match-d=.terraform/* .
```

```
--------------------------------------------------------------------------------
Language                      files          blank        comment           code
--------------------------------------------------------------------------------
Text                              2              0              0          10127
Go                               46            553            250           2628
YAML                              7             25            106            855
Markdown                         13            207              0            347
HTML                              1             43              0            254
JSON                              1              0              0            199
make                              8             41              1            129
Python                            3             18              0             83
Dockerfile                        6             26             15             61
Bourne Shell                      5             10              0             30
SQL                               2              8             10             25
Bourne Again Shell                1              2              3              6
--------------------------------------------------------------------------------
SUM:                             95            933            385          14744
--------------------------------------------------------------------------------
```

### local load testing

Using k6, I ran queries against random databases with random queries of length 1-4. This is with the script running locally against localhost, where network effects are minimal. Ramping up to 400 VUs, we see a throughput of 222 rps (if I'm interpretaing the K6 output correctly). I think the `serve` component can do better. I also think 


```
     ✗ status was 200
      ↳  91% — ✓ 10331 / ✗ 996

     checks.........................: 91.20% 10331 out of 11327
     data_received..................: 8.8 MB 172 kB/s
     data_sent......................: 1.6 MB 32 kB/s
     http_req_blocked...............: avg=64.22µs min=0s      med=7.73µs   max=17.45ms  p(90)=224.12µs p(95)=366.5µs 
     http_req_connecting............: avg=41.35µs min=0s      med=0s       max=17.39ms  p(90)=149.31µs p(95)=244.1µs 
   ✓ http_req_duration..............: avg=44.53ms min=0s      med=1.09ms   max=1.12s    p(90)=14.91ms  p(95)=33.88ms 
       { expected_response:true }...: avg=48.74ms min=113.6µs med=1.2ms    max=1.12s    p(90)=16.28ms  p(95)=39.18ms 
     http_req_failed................: 8.79%  996 out of 11327
     http_req_receiving.............: avg=60.09µs min=0s      med=49.16µs  max=833.26µs p(90)=119.57µs p(95)=139.16µs
     http_req_sending...............: avg=32.07µs min=0s      med=25.43µs  max=1.92ms   p(90)=60.44µs  p(95)=86.18µs 
     http_req_tls_handshaking.......: avg=0s      min=0s      med=0s       max=0s       p(90)=0s       p(95)=0s      
     http_req_waiting...............: avg=44.44ms min=0s      med=977.91µs max=1.12s    p(90)=14.82ms  p(95)=33.82ms 
     http_reqs......................: 11327  222.08811/s
     iteration_duration.............: avg=1.04s   min=1s      med=1s       max=2.12s    p(90)=1.01s    p(95)=1.03s   
     iterations.....................: 11327  222.08811/s
     vus............................: 4      min=4              max=399
     vus_max........................: 400    min=400            max=400
```

Points of comparison:

* https://stackoverflow.com/questions/373098/whats-the-average-requests-per-second-for-a-production-web-application
* https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/
* https://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/