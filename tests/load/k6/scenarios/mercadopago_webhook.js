import http from 'k6/http';
import { check } from 'k6';
import { BASE_URL, MP_WEBHOOK_TOKEN, defaultWebhookOptions } from '../lib/config.js';
import { buildMarketPagoWebhook } from '../lib/builders.js';

export const options = defaultWebhookOptions();

export default function () {
  const payload = buildMarketPagoWebhook(100000 + __VU + __ITER);
  const headers = {
    'Content-Type': 'application/json',
  };

  if (MP_WEBHOOK_TOKEN) {
    headers['X-Webhook-Token'] = MP_WEBHOOK_TOKEN;
  }

  const res = http.post(`${BASE_URL}/webhook/mercadopago`, payload, { headers });

  check(res, {
    'mercado pago webhook returned success': (r) => r.status === 200 || r.status === 401,
  });
}
