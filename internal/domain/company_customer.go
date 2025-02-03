package domain

import (
	"context"
)

type CompanyCustomer struct {
	ID            int64   `json:"id"`
	Company       string  `json:"company"`
	UserID        int64   `json:"user_id"`
	FullName      string  `json:"full_name"`
	Phone         string  `json:"phone"`
	Turnover      float64 `json:"turnover"`
	LastVisitDate string  `json:"last_visit_date"`
	VisitsCount   int64   `json:"visits_count"`
	AverageBill   float64 `json:"average_bill"`
}

type CompanyCustomerRepository interface {
	Store(ctx context.Context, cc *CompanyCustomer) error
	ExistsByCompanyUserId(ctx context.Context, company string, userId int64) (bool, error)
	UpdateByCompanyUserId(ctx context.Context, cc *CompanyCustomer) error
}
