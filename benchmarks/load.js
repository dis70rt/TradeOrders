import http from 'k6/http';
import { uuidv4 } from 'https://jslib.k6.io/k6-utils/1.4.0/index.js';

const symbols = [
  "BTC-USD",
  "ETH-USD",
  "SOL-USD",
  "BNB-USD",
  "XRP-USD",
  "ADA-USD",
  "DOGE-USD",
  "DOT-USD",
  "MATIC-USD",
  "LTC-USD",
];

export const options = {
  scenarios: {
    find_limit: {
      executor: 'ramping-arrival-rate',
      timeUnit: '1s',
      preAllocatedVUs: 5000,
      maxVUs: 10000,
      stages: [
        { target: 12000, duration: '20s' },
        { target: 20000, duration: '20s' },
        { target: 30000, duration: '20s' },
        { target: 40000, duration: '20s' },
      ],
      gracefulStop: '30s',
    },
  },
};

export default function () {
  const instrument = symbols[Math.floor(Math.random() * symbols.length)];
  let q = Number((0.1 + Math.random() * 2).toFixed(4));
  let p = Number((69000 + Math.random() * 3000).toFixed(2));

  let body = JSON.stringify({
    client_id: uuidv4(),
    instrument: instrument,
    side: Math.random() > 0.5 ? 'BUY' : 'SELL',
    type: 'LIMIT',
    price: p,
    quantity: q,
    remaining: q,
  });

  http.post('http://localhost:8080/api/v1/orders/', body, {
    headers: {
      'Content-Type': 'application/json',
      'Idempotency-Key': uuidv4(),
    },
  });
}
