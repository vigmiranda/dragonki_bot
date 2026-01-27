package postgres

import (
	"database/sql"
	"time"

	"github.com/vigmiranda/dragonki/internal/domain"
)

type PaymentRepo struct {
	db *sql.DB
}

func NewPaymentRepo(postgres *Postgres) *PaymentRepo {
	return &PaymentRepo{db: postgres.DB}
}

func (r *PaymentRepo) Save(p *domain.Payment) error {
	_, err := r.db.Exec(`
		INSERT INTO payments (id, player_id, amount, status, qr_code, qr_code_base64, partner_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, p.ID, p.PlayerID, p.Amount, p.Status, p.QRCode, p.QRCodeBase64, p.PartnerId, time.Now())
	return err
}

func (r *PaymentRepo) FindByPartnerID(id int64) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.QueryRow(`
		SELECT id, player_id, amount, status,  qr_code, qr_code_base64, created_at
		FROM payments WHERE partner_id=$1
	`, id).Scan(&p.ID, &p.PlayerID, &p.Amount, &p.Status, &p.QRCode, &p.QRCodeBase64, &p.CreatedAt)
	return &p, err
}

func (r *PaymentRepo) UpdateStatus(id string, status domain.PaymentStatus) error {
	_, err := r.db.Exec(`
		UPDATE payments SET status=$1 WHERE id=$2
	`, status, id)
	return err
}
