package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/choco_parser_go/internal/domain"
	trmpgx "github.com/ibookerke/choco_parser_go/internal/pkg/pgx"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type CompanyCustomerRepository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	trm    trm.Manager
}

func NewCompanyCustomerRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trm trm.Manager,
) *CompanyCustomerRepository {
	return &CompanyCustomerRepository{
		pool:   pool,
		getter: getter,
		trm:    trm,
	}
}

const (
	companyCustomerInsertSQL = `INSERT INTO company_customers
    (company, user_id, full_name, phone, turnover, last_visit_date, visits_count, average_bill)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	companyCustomerExistsSQL = `SELECT EXISTS (
		SELECT 1 
		FROM company_customers 
		WHERE company = $1 AND user_id = $2 LIMIT 1)`

	companyCustomerUpdateSQL = `UPDATE company_customers
    SET full_name = $1, phone = $2, turnover = $3, last_visit_date = $4, visits_count = $5, average_bill = $6
    WHERE company = $7 AND user_id = $8`
)

func (r *CompanyCustomerRepository) Store(ctx context.Context, cc *domain.CompanyCustomer) error {
	exec := r.getter.DefaultTrOrDB(ctx, r.pool)
	_, err := exec.Exec(ctx, companyCustomerInsertSQL,
		cc.Company,
		cc.UserID,
		cc.FullName,
		cc.Phone,
		cc.Turnover,
		cc.LastVisitDate,
		cc.VisitsCount,
		cc.AverageBill,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *CompanyCustomerRepository) ExistsByCompanyUserId(ctx context.Context, company string, userId int64) (bool, error) {
	exec := r.getter.DefaultTrOrDB(ctx, r.pool)
	var exists bool
	err := exec.QueryRow(ctx, companyCustomerExistsSQL, company, userId).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *CompanyCustomerRepository) UpdateByCompanyUserId(ctx context.Context, cc *domain.CompanyCustomer) error {
	exec := r.getter.DefaultTrOrDB(ctx, r.pool)
	_, err := exec.Exec(ctx, companyCustomerUpdateSQL,
		cc.FullName,
		cc.Phone,
		cc.Turnover,
		cc.LastVisitDate,
		cc.VisitsCount,
		cc.AverageBill,
		cc.Company,
		cc.UserID,
	)
	if err != nil {
		return err
	}
	return nil
}
