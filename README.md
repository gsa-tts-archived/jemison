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
> 1. **MAINTAINABLE**: 20x fewer lines of code
> 2. **AFFORDABLE** 10x less resource intenstive (read: hosting costs)
> 3. **BLAZING FAST**: 10x more performant (and scales horitzontally)

Why `jemison` is an MVP candidate for search:

1. `jemison` can already crawl and index HTML and PDFs, handling diacritics (e.g. Español).
2. `jemison` can deploy a full crawl/index/serve solution in 1GB of RAM on [cloud.gov](https://cloud.gov/).
3. `jemison` is 2500 lines of code: easy to maintain, with an extensible, distributed architecture.
4. `jemison` can handle 55K requests per minute (10x the current search.gov load) with 256**M**B of RAM.

I lied. There is no #5.


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
cloc --exclude-ext=yml,yaml,html,css,less,js,json,svg,scss --fullpath --not-match-d=terraform/\.terraform .
```

```
--------------------------------------------------------------------------------
Language                      files          blank        comment           code
--------------------------------------------------------------------------------
Go                               53            639            320           2985
JSON                              2              0              0           2852
Markdown                         17            257              0            451
HCL                               5             48             37            291
make                              9             60              3            240
Python                            3             18              0             82
Dockerfile                        6             27             15             64
Bourne Shell                      5             10              0             30
SQL                               2              8             11             26
Bourne Again Shell                2              3              5             11
--------------------------------------------------------------------------------
SUM:                            104           1070            391           7032
--------------------------------------------------------------------------------
```

### local load testing

Using k6, I ran queries against random databases with random queries of length 1-4. This is with the script running locally against localhost, where network effects are minimal. Ramping up to 1000 simultaneous users, we see a throughput of 930 requests per second. 

That means that the search server can handle 55K (fifty-five thousand) requests per minute. The server in question **only has 256M of memory**. Or, if we wanted to scale to 200K requests per minute, we could do so by paying for 1GB of RAM/month, which has an annual cost of $130.

search.gov sees approximately 5K requests per minute; with a tiny fraction of the resources of the current infrastructre, this design can handle 10x the load without breaking a sweat.


```
         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: script2.js
        output: -

     scenarios: (100.00%) 1 scenario, 1200 max VUs, 1m0s max duration (incl. graceful stop):
              * constant_request_rate: 1000.00 iterations/s for 30s (maxVUs: 100-1200, gracefulStop: 30s)


     ✓ status was 200

     checks.........................: 100.00% 28821 out of 28821
     data_received..................: 25 MB   811 kB/s
     data_sent......................: 4.5 MB  146 kB/s
     dropped_iterations.............: 1180    38.051839/s
     http_req_blocked...............: avg=12.03µs min=868ns    med=3.41µs  max=11.65ms  p(90)=7.7µs   p(95)=12.27µs
     http_req_connecting............: avg=5.89µs  min=0s       med=0s      max=11.61ms  p(90)=0s      p(95)=0s     
     http_req_duration..............: avg=5.65ms  min=285.16µs med=2.32ms  max=130.44ms p(90)=13.86ms p(95)=22.32ms
       { expected_response:true }...: avg=5.65ms  min=285.16µs med=2.32ms  max=130.44ms p(90)=13.86ms p(95)=22.32ms
     http_req_failed................: 0.00%   0 out of 28821
     http_req_receiving.............: avg=33.72µs min=4.28µs   med=26.59µs max=3.18ms   p(90)=56.84µs p(95)=74.05µs
     http_req_sending...............: avg=15.69µs min=2.6µs    med=10.73µs max=2.08ms   p(90)=26.83µs p(95)=36.41µs
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s       p(90)=0s      p(95)=0s     
     http_req_waiting...............: avg=5.6ms   min=272.58µs med=2.26ms  max=130.36ms p(90)=13.78ms p(95)=22.23ms
     http_reqs......................: 28821   929.400047/s
     iteration_duration.............: avg=1s      min=1s       med=1s      max=1.13s    p(90)=1.01s   p(95)=1.02s  
     iterations.....................: 28821   929.400047/s
     vus............................: 11      min=11             max=1032
     vus_max........................: 1041    min=485            max=1041
```

Points of comparison:

* https://stackoverflow.com/questions/373098/whats-the-average-requests-per-second-for-a-production-web-application
* https://www.freecodecamp.org/news/million-websockets-and-go-cc58418460bb/
* https://tleyden.github.io/blog/2016/11/21/tuning-the-go-http-client-library-for-load-testing/