export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';
export const TELEGRAM_CHAT_ID = Number(__ENV.TELEGRAM_CHAT_ID || 200);
export const TELEGRAM_USER_ID = Number(__ENV.TELEGRAM_USER_ID || 200);
export const TELEGRAM_MESSAGE_ID = Number(__ENV.TELEGRAM_MESSAGE_ID || 1000);
export const MONSTER_ID = Number(__ENV.MONSTER_ID || 1);
export const MP_WEBHOOK_TOKEN = __ENV.MP_WEBHOOK_TOKEN || '';

export function defaultTelegramOptions() {
  return {
    stages: [
      { duration: '30s', target: 100 },
      { duration: '1m', target: 100 },
      { duration: '30s', target: 0 },
    ],
    thresholds: {
      http_req_failed: ['rate<0.01'],
      http_req_duration: ['p(95)<300', 'p(99)<700'],
    },
  };
}

export function defaultWebhookOptions() {
  return {
    stages: [
      { duration: '20s', target: 50 },
      { duration: '40s', target: 50 },
      { duration: '20s', target: 0 },
    ],
    thresholds: {
      http_req_failed: ['rate<0.01'],
      http_req_duration: ['p(95)<250', 'p(99)<500'],
    },
  };
}
