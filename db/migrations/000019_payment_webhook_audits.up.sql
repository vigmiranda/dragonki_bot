CREATE TABLE IF NOT EXISTS payment_webhook_audits (
    id BIGSERIAL PRIMARY KEY,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    request_method TEXT NOT NULL,
    request_path TEXT NOT NULL,
    remote_addr TEXT NOT NULL DEFAULT '',
    headers_json TEXT NOT NULL DEFAULT '{}',
    query_params_json TEXT NOT NULL DEFAULT '{}',
    raw_body TEXT NOT NULL DEFAULT '',
    event_type TEXT NOT NULL DEFAULT '',
    event_action TEXT NOT NULL DEFAULT '',
    event_data_id TEXT NOT NULL DEFAULT '',
    partner_payment_id BIGINT NOT NULL DEFAULT 0,
    authorized BOOLEAN NOT NULL DEFAULT FALSE,
    feature_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    processing_status TEXT NOT NULL,
    http_status INT NOT NULL,
    error_message TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payment_webhook_audits_received_at
    ON payment_webhook_audits (received_at DESC);

CREATE INDEX IF NOT EXISTS idx_payment_webhook_audits_partner_payment_id
    ON payment_webhook_audits (partner_payment_id);

CREATE INDEX IF NOT EXISTS idx_payment_webhook_audits_processing_status
    ON payment_webhook_audits (processing_status);
