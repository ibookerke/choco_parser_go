package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/choco_parser_go/internal/domain"
	trmpgx "github.com/ibookerke/choco_parser_go/internal/pkg/pgx"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type CustomerRepository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	trm    trm.Manager
}

func NewCustomerRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trm trm.Manager,
) *CustomerRepository {
	return &CustomerRepository{
		pool:   pool,
		getter: getter,
		trm:    trm,
	}
}

const (
	customerExistsById = `SELECT 
		EXISTS ( SELECT 1 
			FROM customers 
			WHERE id = $1 LIMIT 1)`

	customerCreateSql = `INSERT INTO customers
		(id, user_id, phone, birthday, full_name, orders_count) 
	VALUES 
		($1, $2, $3, $4, $5, $6)`

	findCustomerById = `SELECT
 		id, user_id, phone, birthday, full_name, orders_count
	FROM customers
	WHERE id = $1`
)

func (c *CustomerRepository) ExistsById(ctx context.Context, id domain.CustomerID) (bool, error) {
	exec := c.getter.DefaultTrOrDB(ctx, c.pool)

	var exists bool
	err := exec.QueryRow(ctx, customerExistsById, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (c *CustomerRepository) Create(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	exec := c.getter.DefaultTrOrDB(ctx, c.pool)

	_, err := exec.Exec(
		ctx,
		customerCreateSql,
		customer.ID,
		customer.UserID,
		customer.Phone,
		customer.Birthday,
		customer.FullName,
		customer.OrderCount,
	)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (c *CustomerRepository) FindById(ctx context.Context, id domain.CustomerID) (domain.Customer, error) {
	exec := c.getter.DefaultTrOrDB(ctx, c.pool)

	var customer domain.Customer
	err := exec.QueryRow(ctx, findCustomerById, id).Scan(
		&customer.ID,
		&customer.UserID,
		&customer.Phone,
		&customer.Birthday,
		&customer.FullName,
		&customer.OrderCount,
	)
	if err != nil {
		return domain.Customer{}, err
	}

	return customer, nil
}
