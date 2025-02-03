package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/choco_parser_go/internal/domain"
	trmpgx "github.com/ibookerke/choco_parser_go/internal/pkg/pgx"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type PaymentRepository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	trm    trm.Manager
}

func NewPaymentRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trm trm.Manager,
) *PaymentRepository {
	return &PaymentRepository{
		pool:   pool,
		getter: getter,
		trm:    trm,
	}
}

const (
	paymentCreateSql = `INSERT INTO payments
    (id, created_by, type, amount, discount_amount, created_at, location_title, location_partner_id)
    values ($1, $2, $3, $4, $5, $6, $7, $8)`

	paymentExistsById = `SELECT 
		EXISTS ( SELECT 1 
			FROM payments 
			WHERE id = $1 LIMIT 1)`

	paymentFindById = `SELECT
		id, created_by, type, amount, discount_amount, created_at, location_title, location_partner_id
	FROM payments
	WHERE id = $1`
)

func (p *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	exec := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := exec.Exec(ctx, paymentCreateSql,
		payment.ID,
		payment.CreatedBy,
		payment.Type,
		payment.Amount,
		payment.DiscountAmount,
		payment.CreatedAt,
		payment.LocationTitle,
		payment.LocationPartnerID,
	)
	if err != nil {
		return nil, err
	}
	return payment, nil

}

func (p *PaymentRepository) ExistsById(ctx context.Context, id domain.PaymentID) (bool, error) {
	exec := p.getter.DefaultTrOrDB(ctx, p.pool)
	var exists bool
	err := exec.QueryRow(ctx, paymentExistsById, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (p *PaymentRepository) FindById(ctx context.Context, id domain.PaymentID) (*domain.Payment, error) {
	exec := p.getter.DefaultTrOrDB(ctx, p.pool)
	var payment domain.Payment
	err := exec.QueryRow(ctx, paymentFindById, id).Scan(
		&payment.ID,
		&payment.CreatedBy,
		&payment.Type,
		&payment.Amount,
		&payment.DiscountAmount,
		&payment.CreatedAt,
		&payment.LocationTitle,
		&payment.LocationPartnerID,
	)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
