package domain

import (
	"context"
	"strconv"
)

type CustomerID int64

func CustomerIdToStr(id CustomerID) string {
	return strconv.FormatInt(int64(id), 10)
}

type Customer struct {
	ID         CustomerID `json:"id"`
	UserID     int        `json:"user_id"`
	Phone      string     `json:"phone,omitempty"`
	Birthday   string     `json:"birthday,omitempty"`
	FullName   string     `json:"full_name,omitempty"`
	OrderCount int64      `json:"orderCount,omitempty"`
}

type CustomerRepository interface {
	ExistsById(ctx context.Context, id CustomerID) (bool, error)
	Create(ctx context.Context, customer *Customer) (*Customer, error)
	FindById(ctx context.Context, id CustomerID) (Customer, error)
}
