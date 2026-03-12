import http from 'k6/http';
import { check } from 'k6';
import { BASE_URL, defaultTelegramOptions } from '../lib/config.js';
import { buildBigBossCallback } from '../lib/builders.js';

export const options = {
  ...defaultTelegramOptions(),
  stages: [
    { duration: '10s', target: 0 },
    { duration: '10s', target: 200 },
    { duration: '20s', target: 200 },
    { duration: '10s', target: 0 },
  ],
};

export default function () {
  const payload = JSON.stringify(buildBigBossCallback());
  const res = http.post(`${BASE_URL}/webhook/telegram`, payload, {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  check(res, {
    'telegram big boss returned 200': (r) => r.status === 200,
  });
}
