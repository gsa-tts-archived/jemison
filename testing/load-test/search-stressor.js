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
      rate: 500, 
      timeUnit: '1s', 
      duration: '30s',
      preAllocatedVUs: 100,
      maxVUs: 1000, 
      startTime: "40s"
    },
    teeter_totter: {
      executor: 'ramping-vus',
      startVUs: 10,
      gracefulRampDown: "10s",
      stages: [
        { duration: "10s", target: 100 },
        { duration: "15s", target: 250 },
        { duration: "15s", target: 750 },
        { duration: "10s", target: 100 },
      ],
      startTime: "70s"
    },
  },
  
};

const wordlist = JSON.parse(open('./wordlist.json')); 

// Simulated user behavior
export default function () {
  let terms = "";
  let number_of_terms = Math.floor(Math.random() * 4) + 1;

  for (let i = 0; i < number_of_terms ; i++) {
    let choice = Math.floor(Math.random() * wordlist.length)
    let word = wordlist[choice];
    terms += word + " ";  
  }

  let data = { 
    host: "gov", 
    d64_start: "0", 
    d64_end: Number("0x01FF000000000000").toString(), 
    terms: terms.trim()};

  let res = http.post(
    "http://localhost:10000/api/search",
    JSON.stringify(data)
  )

  // Validate response status
  check(res, { "status was 200": (r) => r.status == 200 });
  sleep(1);
}