CREATE TABLE IF NOT EXISTS telegram_callback_receipts (
    callback_id TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    message_id BIGINT NOT NULL,
    action TEXT NOT NULL,
    raw_data TEXT NOT NULL,
    received_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_telegram_callback_receipts_received_at
    ON telegram_callback_receipts (received_at DESC);

CREATE INDEX IF NOT EXISTS idx_telegram_callback_receipts_action
    ON telegram_callback_receipts (action);
