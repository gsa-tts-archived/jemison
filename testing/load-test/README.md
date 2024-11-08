# load testing

This is not containerized/automated at this time.

1. You will need to install `k6` (https://k6.io/)
2. You will need Python

## generate wordlists

The tests pull random words from existing databases to run authentic, random queries.

Edit `get_distinct_words.py` to include names of SQLite databases that have been generated locally. For example, if you have crawled `search.gov`, or have some other set of test databases, list those here.

Then, run

```
python get_distinct_words.py
```

This will generate the file `wordlists.json`. The structure of this file is:

```
{
  "site.gov": [
    "word",
    "anotherword",
    "somanywords"
  ],
  "another.gov": [ ... ]
}
```

## run k6

These stress tests work with a `serve` component with at least 512M of RAM. Less RAM will likely crash the container because the RAM limit is hit. Golang scales connection handling very well with RAM. That is, as it works to handle hundreds or thousands of concurrent connections, it will consume RAM. With 512M of RAM, we can handle around 1000 simultaneous connections for a throughput of ~55K RPM (requests per minute).

To run the test suite

```
k6 run search-stressor.js
```

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