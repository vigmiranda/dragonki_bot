import { TELEGRAM_CHAT_ID, TELEGRAM_MESSAGE_ID, TELEGRAM_USER_ID } from './config.js';

function telegramIdentity() {
  const userID = TELEGRAM_USER_ID + __VU;
  const chatID = TELEGRAM_CHAT_ID + __VU;
  const messageID = TELEGRAM_MESSAGE_ID + (__VU * 100000) + __ITER;

  return { userID, chatID, messageID };
}

function callbackQuery(action, idSuffix, extra = {}) {
  const identity = telegramIdentity();

  return {
    update_id: Date.now() + idSuffix,
    callback_query: {
      id: `cb-${idSuffix}-${__VU}-${__ITER}`,
      from: {
        id: identity.userID,
        is_bot: false,
        first_name: 'Stress',
        username: `stress_${__VU}`,
      },
      message: {
        message_id: identity.messageID,
        chat: {
          id: identity.chatID,
          type: 'private',
        },
      },
      data: action,
      ...extra,
    },
  };
}

export function buildBattleCallback(monsterID) {
  return callbackQuery(`action_fight?monsterID=${monsterID}`, 1);
}

export function buildBigBossCallback() {
  return callbackQuery('big_boss_attack', 2);
}

export function buildMarketPagoWebhook(partnerPaymentID) {
  return JSON.stringify({
    type: 'payment',
    action: 'payment.updated',
    data: {
      id: String(partnerPaymentID),
    },
  });
}
