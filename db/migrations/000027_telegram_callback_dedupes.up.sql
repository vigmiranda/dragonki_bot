CREATE TABLE IF NOT EXISTS telegram_callback_dedupes (
    dedupe_key TEXT PRIMARY KEY,
    action TEXT NOT NULL,
    window_start TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_telegram_callback_dedupes_window_start
    ON telegram_callback_dedupes (window_start DESC);

CREATE INDEX IF NOT EXISTS idx_telegram_callback_dedupes_action
    ON telegram_callback_dedupes (action);
