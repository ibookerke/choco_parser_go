package domain

import "context"

type Payment struct {
	ID                PaymentID `json:"id"`
	Type              string    `json:"type"`
	CreatedBy         int64     `json:"created_by,omitempty"`
	Amount            int64     `json:"amount,omitempty"`
	DiscountAmount    int64     `json:"discount_amount,omitempty"`
	CreatedAt         string    `json:"created_at,omitempty"`
	LocationTitle     string    `json:"location_title,omitempty"`
	LocationPartnerID string    `json:"location_partner_id,omitempty"`
}

type PaymentID int64

type PaymentRepository interface {
	ExistsById(ctx context.Context, id PaymentID) (bool, error)
	FindById(ctx context.Context, id PaymentID) (*Payment, error)
	Create(ctx context.Context, payment *Payment) (*Payment, error)
}
