import http from "k6/http";
import { check, sleep } from "k6";
import { Rate, Trend } from "k6/metrics";

// Custom metrics
export const errorRate = new Rate("errors");
export const responseTime = new Trend("response_time");

export const options = {
  stages: [
    { duration: "10s", target: 10 }, // Warm-up phase
    { duration: "30s", target: 50 }, // Steady ramp-up
    { duration: "60s", target: 100 }, // Sustained high load
    { duration: "20s", target: 50 }, // Gradual ramp-down
    { duration: "10s", target: 10 }, // Cool-down phase
  ],
  thresholds: {
    errors: ["rate<0.01"], // Less than 1% errors
    http_req_duration: ["p(95)<500"], // 95% of requests should be under 500ms
  },
};

export default function () {
  const res = http.get("https://trendingrepos.glup3.dev");

  // Validate response
  const success = check(res, {
    "status was 200": (r) => r.status === 200,
    "response time < 500ms": (r) => r.timings.duration < 500,
  });

  // Track metrics
  errorRate.add(!success);
  responseTime.add(res.timings.duration);

  sleep(1); // Simulate user think time
}
