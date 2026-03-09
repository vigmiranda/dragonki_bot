package domain

import "time"

type PaymentWebhookAudit struct {
	ID               uint      `gorm:"primaryKey"`
	ReceivedAt       time.Time `gorm:"column:received_at;not null"`
	RequestMethod    string    `gorm:"column:request_method;not null"`
	RequestPath      string    `gorm:"column:request_path;not null"`
	RemoteAddr       string    `gorm:"column:remote_addr"`
	HeadersJSON      string    `gorm:"column:headers_json;not null"`
	QueryParamsJSON  string    `gorm:"column:query_params_json;not null"`
	RawBody          string    `gorm:"column:raw_body;not null"`
	EventType        string    `gorm:"column:event_type;not null"`
	EventAction      string    `gorm:"column:event_action;not null"`
	EventDataID      string    `gorm:"column:event_data_id;not null"`
	PartnerPaymentID int64     `gorm:"column:partner_payment_id"`
	Authorized       bool      `gorm:"column:authorized;not null"`
	FeatureEnabled   bool      `gorm:"column:feature_enabled;not null"`
	ProcessingStatus string    `gorm:"column:processing_status;not null"`
	HTTPStatus       int       `gorm:"column:http_status;not null"`
	ErrorMessage     string    `gorm:"column:error_message"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (PaymentWebhookAudit) TableName() string {
	return "payment_webhook_audits"
}
