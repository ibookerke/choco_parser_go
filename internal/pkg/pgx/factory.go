package pgx

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

// NewDefaultFactory creates default trm.Transaction(pgx.Tx).
func NewDefaultFactory(db *pgxpool.Pool) trm.TrFactory {
	return func(ctx context.Context, trms trm.Settings) (context.Context, trm.Transaction, error) {
		s, _ := trms.(Settings)

		return NewTransaction(ctx, s.TxOpts(), db)
	}
}
