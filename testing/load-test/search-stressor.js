import http from "k6/http";
import { check, sleep } from "k6";

// https://k6.io/blog/how-to-generate-a-constant-request-rate-with-the-new-scenarios-api/
export const options = {
  thresholds: {
    // http errors should be less than 1%
    http_req_failed: ['rate<0.01'], 
    // 95% of requests should be below 2000ms
    http_req_duration: ['p(95)<2000'], 
  },
  scenarios: {
    one_is_such_a_lonely_number: {
      executor: 'constant-arrival-rate',
      rate: 1, 
      timeUnit: '1s', 
      duration: '10s',
      preAllocatedVUs: 10, 
      maxVUs: 100, 
    },
    the_marching_100: {
      executor: 'constant-arrival-rate',
      rate: 100, 
      timeUnit: '1s', 
      duration: '30s',
      preAllocatedVUs: 100, 
      maxVUs: 200, 
      startTime: "10s"
    },
    a_new_millenium: {
      executor: 'constant-arrival-rate',
      rate: 1000, 
      timeUnit: '1s', 
      duration: '60s',
      preAllocatedVUs: 100,
      maxVUs: 1500, 
      startTime: "40s"
    },
    teeter_totter: {
      executor: 'ramping-vus',
      startVUs: 10,
      gracefulRampDown: "10s",
      stages: [
        { duration: "10s", target: 100 },
        { duration: "20s", target: 500 },
        { duration: "60s", target: 1200 },
        { duration: "30s", target: 100 },
      ],
      startTime: "100s"
    },
  },
  
};

const wordlists = JSON.parse(open('./wordlists.json')); 

// Simulated user behavior
export default function () {
  let hosts = Object.keys(wordlists);
  let terms = "";
  let number_of_terms = Math.floor(Math.random() * 3) + 1;
  let which_host_num =  Math.floor(Math.random() * hosts.length);
  let which_host = hosts[which_host_num]

  for (let i = 0; i < number_of_terms ; i++) {
    var low = wordlists[which_host]
    if (low && low.length > 0) {
      let choice = Math.floor(Math.random() * low.length)
      let word = low[choice];
      terms += " " + word;  
    }
  }

  let data = { host: which_host, terms: terms};

  let res = http.post(
    "http://localhost:10000/api/search",
    JSON.stringify(data)
  )

  // Validate response status
  check(res, { "status was 200": (r) => r.status == 200 });
  sleep(1);
}