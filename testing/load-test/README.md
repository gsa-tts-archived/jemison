# load testing

This is not containerized/automated at this time.

1. You will need to install `k6` (https://k6.io/)
2. You will need Python

## generate wordlists

The tests pull random words from existing database to run authentic, random queries.

Run

```
python get_distinct_words.py
```

This will generate/overwrite the file `wordlist.json`. The structure of this file is:

```
[
    "word",
    "anotherword",
    "somanywords"
]
```

## run k6

These stress tests work with a `serve` component with at least 512M of RAM. Less RAM will likely crash the container because the RAM limit is hit. Golang scales connection handling very well with RAM. That is, as it works to handle hundreds or thousands of concurrent connections, it will consume RAM. With 512M of RAM, we can handle around 500 simultaneous connections for a throughput of ~12K RPM (requests per minute).

To run the test suite

```
k6 run search-stressor.js
```


## Postgres results

With a connection pool limit of 100, and no particular optimizations on the API (running in `debug` mode), we see roughly 200 req/s. This comes out to around 12K req/minute, which is twice what Search.gov sees around peak (5K req/minute).

We can make this work with a single Postgres instance running FTS and a single Golang API server.

These numbers are against a local stack, and therefore have no network overhead. The question is less about the transmission time over the wire and more about what the upper bounds are on load, which a local test works well for.

More testing is needed to see if there are failure modes here. It seems, from the test, that we returned 200 to everything, but some requests pushed out to 4s. (Again, this is my inexpert reading of K6 output.)


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


























## Sqlite performance

When we were backing things with SQlite, we saw numbers up to around 600 req/sec.

You should get output similar to

```
load-test$ k6 run search-stressor.js 

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


     ✓ status was 200

     checks.........................: 100.00% 135962 out of 135962
     data_received..................: 96 MB   434 kB/s
     data_sent......................: 22 MB   100 kB/s
     dropped_iterations.............: 1415    6.403672/s
     http_req_blocked...............: avg=11.5µs  min=851ns    med=4.13µs  max=17.18ms  p(90)=11.48µs p(95)=13.77µs 
     http_req_connecting............: avg=4.1µs   min=0s       med=0s      max=17.12ms  p(90)=0s      p(95)=0s      
   ✓ http_req_duration..............: avg=26.68ms min=310.41µs med=9.84ms  max=513.96ms p(90)=77.61ms p(95)=94.25ms 
       { expected_response:true }...: avg=26.68ms min=310.41µs med=9.84ms  max=513.96ms p(90)=77.61ms p(95)=94.25ms 
   ✓ http_req_failed................: 0.00%   0 out of 135962
     http_req_receiving.............: avg=48.08µs min=4.55µs   med=28.53µs max=7.72ms   p(90)=94.41µs p(95)=130.04µs
     http_req_sending...............: avg=23.29µs min=2.63µs   med=13.64µs max=10.57ms  p(90)=43.1µs  p(95)=54.39µs 
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s      max=0s       p(90)=0s      p(95)=0s      
     http_req_waiting...............: avg=26.61ms min=285.86µs med=9.74ms  max=513.43ms p(90)=77.54ms p(95)=94.18ms 
     http_reqs......................: 135962  615.30462/s
     iteration_duration.............: avg=1.02s   min=1s       med=1.01s   max=1.52s    p(90)=1.07s   p(95)=1.09s   
     iterations.....................: 135962  615.30462/s
     vus............................: 9       min=1                max=1198
     vus_max........................: 1999    min=1200             max=1999


running (3m41.0s), 0000/1999 VUs, 135962 complete and 0 interrupted iterations
one_is_such_a_lonely_number ✓ [======================================] 000/010 VUs    10s   1.00 iters/s
the_marching_100            ✓ [======================================] 000/103 VUs    30s   100.00 iters/s
a_new_millenium             ✓ [======================================] 0000/1047 VUs  1m0s  1000.00 iters/s
teeter_totter               ✓ [======================================] 0000/1200 VUs  2m0s 
```

The test suite might change over time, but as of this writing (2024-11-08) this is representative.