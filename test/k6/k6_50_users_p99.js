import http from "k6/http";
import { sleep } from "k6";
import { templateSearch } from "./templateSearch.js";

export const options = {
  vus: 51,
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

  console.log(
    `search => ${item}\nresult=>${response.message} ${data.length} ${
      data.length !== 0 ? data[0].product_name : "no product available"
    }`
  );

  // sleep(Math.random() * 2);
}

// export function handleSummary(data) {
//   return {
//     "result_k6_50_users_p99.json": JSON.stringify(data),
//   };
// }
