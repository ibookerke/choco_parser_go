package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
	trmcontext "github.com/ibookerke/choco_parser_go/internal/pkg/trm/context"
)

// DefaultCtxGetter is the CtxGetter with settings.DefaultCtxKey.
//
//nolint:gochecknoglobals
var DefaultCtxGetter = NewCtxGetter(trmcontext.DefaultManager)

// CtxGetter gets Executor from trm.СtxManager by casting trm.Transaction to Tr.
type CtxGetter struct {
	ctxManager trm.СtxManager
}

//revive:disable:exported
func NewCtxGetter(c trm.СtxManager) *CtxGetter {
	return &CtxGetter{ctxManager: c}
}

func (c *CtxGetter) DefaultTrOrDB(ctx context.Context, db Executor) Executor {
	if tr := c.ctxManager.Default(ctx); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) TrOrDB(ctx context.Context, key trm.CtxKey, db Executor) Executor {
	if tr := c.ctxManager.ByKey(ctx, key); tr != nil {
		return c.convert(tr)
	}

	return db
}

func (c *CtxGetter) convert(tr trm.Transaction) Executor {
	if tx, ok := tr.Transaction().(pgx.Tx); ok {
		return tx
	}

	return nil
}
