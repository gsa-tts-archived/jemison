import http from "k6/http";
import { check, sleep } from "k6";

// Test configuration
export const options = {
  thresholds: {
    // Assert that 99% of requests finish within 3000ms.
    http_req_duration: ["p(99) < 3000"],
  },
  // Ramp the number of virtual users up and down
  stages: [
    { duration: "15s", target: 200 },
    { duration: "30s", target: 400 },
    { duration: "5s", target: 50 },
  ],
};

let words = [
  "audit",
  "compliance",
  "2024",
  "2023",
  "omb",
  "treasury",
  "education",
  "health",
  "human",
  "services",
  "green",
  "NSAC",
]

let hosts = [
  "fac.gov",
  "jadud.com",
  "search.gov",
  "cloud.gov",
]

// Simulated user behavior
export default function () {
  let terms = "";
  let number_of_terms = Math.floor(Math.random() * 4);
  let which_host =  Math.floor(Math.random() * hosts.length);

  for (let i = 0; i < number_of_terms ; i++) {
    var word = words[Math.floor(Math.random() * words.length)];
    terms += " " + word;
  }

  let data = { host: hosts[which_host], terms: terms};

  // Test just the stats endpoint
  // let res = http.get(
  //   "http://localhost:10004/api/stats"
  //   // JSON.stringify(data)
  // );

  // let res = http.post(
  //   "http://localhost:10004/api/search",
  //   JSON.stringify(data)
  // )

  let res = http.get("http://localhost:10000/heartbeat")

  // Validate response status
  check(res, { "status was 200": (r) => r.status == 200 });
  sleep(1);
}