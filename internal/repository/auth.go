package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/choco_parser_go/internal/domain"
	trmpgx "github.com/ibookerke/choco_parser_go/internal/pkg/pgx"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

type AuthRepository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	trm    trm.Manager
}

func NewAuthRepository(
	pool *pgxpool.Pool,
	getter *trmpgx.CtxGetter,
	trm trm.Manager,
) *AuthRepository {
	return &AuthRepository{
		pool:   pool,
		getter: getter,
		trm:    trm,
	}
}

const (
	getAuthByClientId = `SELECT 
		id, client_id, token, refresh_token
	FROM auth
		WHERE client_id = $1`

	updateAuthByClientId = `UPDATE auth
		SET token = $1, refresh_token = $2
		WHERE client_id = $3`
)

func (a *AuthRepository) GetAuthByClientId(ctx context.Context, clientId int64) (domain.Auth, error) {
	exec := a.getter.DefaultTrOrDB(ctx, a.pool)
	var auth domain.Auth

	err := exec.QueryRow(
		ctx,
		getAuthByClientId,
		clientId,
	).Scan(
		&auth.ID,
		&auth.ClientID,
		&auth.Token,
		&auth.RefreshToken,
	)
	if err != nil {
		return domain.Auth{}, err
	}
	return auth, nil
}

func (a *AuthRepository) UpdateAuthByClientId(ctx context.Context, token string, refreshToken string, clientId int64) error {
	exec := a.getter.DefaultTrOrDB(ctx, a.pool)
	_, err := exec.Exec(
		ctx,
		updateAuthByClientId,
		token,
		refreshToken,
		clientId,
	)
	if err != nil {
		return err
	}
	return nil
}
