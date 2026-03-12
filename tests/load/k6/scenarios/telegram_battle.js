import http from 'k6/http';
import { check } from 'k6';
import { BASE_URL, MONSTER_ID, defaultTelegramOptions } from '../lib/config.js';
import { buildBattleCallback } from '../lib/builders.js';

export const options = defaultTelegramOptions();

export default function () {
  const payload = JSON.stringify(buildBattleCallback(MONSTER_ID));
  const res = http.post(`${BASE_URL}/webhook/telegram`, payload, {
    headers: {
      'Content-Type': 'application/json',
    },
  });

  check(res, {
    'telegram battle returned 200': (r) => r.status === 200,
  });
}
