package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/choco_parser_go/internal/domain"
	trmpgx "github.com/ibookerke/choco_parser_go/internal/pkg/pgx"

	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type BranchRepository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	trm    trm.Manager
}

func NewBranchRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trm trm.Manager,
) *BranchRepository {
	return &BranchRepository{
		pool:   pool,
		getter: getter,
		trm:    trm,
	}
}

const (
	checkBranchExistById = `SELECT 
		EXISTS ( SELECT 1 
			FROM branches 
			WHERE id = $1 LIMIT 1)`

	branchCreateSql = `INSERT INTO branches
		(id, name, status, type_id, type_name, type_description, token, location_id, location_name, partner_id, partner_name, partner_logo)
    VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	getBranchesByCompanyName = `SELECT id FROM branches 
    WHERE name ILIKE '%' || $1 || '%'`
)

func (b *BranchRepository) CheckIfBranchExist(ctx context.Context, id int64) (bool, error) {
	exec := b.getter.DefaultTrOrDB(ctx, b.pool)
	var exist bool
	err := exec.QueryRow(
		ctx,
		checkBranchExistById,
		id,
	).Scan(
		&exist,
	)
	if err != nil {
		return false, fmt.Errorf("check branch exist: %w", wrapScanError(err))
	}
	return exist, nil
}

func (b *BranchRepository) Create(ctx context.Context, branch *domain.Branch) (*domain.Branch, error) {
	exec := b.getter.DefaultTrOrDB(ctx, b.pool)

	parsedName := strings.ReplaceAll(branch.Name, "\t", " ")

	_, err := exec.Exec(
		ctx,
		branchCreateSql,
		branch.ID,
		parsedName,
		branch.Status,
		branch.TypeID,
		branch.TypeName,
		branch.TypeDescription,
		branch.Token,
		branch.LocationID,
		branch.LocationName,
		branch.PartnerID,
		branch.PartnerName,
		branch.PartnerLogo,
	)
	if err != nil {
		return nil, fmt.Errorf("save branch: %w", wrapScanError(err))
	}

	return branch, nil

}

func (b *BranchRepository) GetAll(ctx context.Context) ([]*domain.Branch, error) {
	exec := b.getter.DefaultTrOrDB(ctx, b.pool)

	rows, err := exec.Query(
		ctx,
		`SELECT 
			id, name, status, type_id, type_name, type_description, token, location_id, location_name, partner_id, partner_name, partner_logo
		FROM branches`,
	)
	if err != nil {
		return nil, fmt.Errorf("get all branches: %w", wrapScanError(err))
	}
	defer rows.Close()

	var branches []*domain.Branch
	for rows.Next() {
		var branch domain.Branch
		err := rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Status,
			&branch.TypeID,
			&branch.TypeName,
			&branch.TypeDescription,
			&branch.Token,
			&branch.LocationID,
			&branch.LocationName,
			&branch.PartnerID,
			&branch.PartnerName,
			&branch.PartnerLogo,
		)
		if err != nil {
			return nil, fmt.Errorf("scan all branches: %w", wrapScanError(err))
		}
		branches = append(branches, &branch)
	}

	return branches, nil
}

func (b *BranchRepository) GetBranchesByCompanyName(ctx context.Context, companyName string) ([]domain.BranchId, error) {
	exec := b.getter.DefaultTrOrDB(ctx, b.pool)

	rows, err := exec.Query(
		ctx,
		getBranchesByCompanyName,
		companyName,
	)
	if err != nil {
		return nil, fmt.Errorf("get branches by company name: %w", wrapScanError(err))
	}
	defer rows.Close()

	var branches []domain.BranchId
	for rows.Next() {
		var branch domain.BranchId
		err := rows.Scan(
			&branch,
		)
		if err != nil {
			return nil, fmt.Errorf("scan branches by company name: %w", wrapScanError(err))
		}
		branches = append(branches, branch)
	}

	return branches, nil
}
