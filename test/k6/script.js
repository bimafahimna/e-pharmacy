import http from "k6/http";
import { sleep } from "k6";
import { templateSearch } from "./templateSearch.js";

export const options = {
  stages: [
    { duration: "1s", target: 5 },
    { duration: "5s", target: 50 },
    { duration: "1s", target: 0 },
  ],
  /* different for each case*/
  vus: 1,
  duration: "300s",
  thresholds: {
    http_req_duration: ["p(99)<500"],
  },
  summaryTrendStats: ["min", "max", "avg", "med", "p(99)"],
};

export default function () {
  const item =
    templateSearch[Math.floor(Math.random() * templateSearch.length)];

  const res = http.get(
    `http://localhost:8000/products?limit=20&page=1&s=${encodeURIComponent(
      item
    )}`
  );
  const response = res.json();
  const data = response.data || [];

  sleep(Math.random() * 2);
}

// export function handleSummary(data) {
//   return {
//     "script_result.json": JSON.stringify(data),
//   };
// }
