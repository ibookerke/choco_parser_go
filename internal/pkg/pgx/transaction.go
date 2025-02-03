package pgx

import (
	"context"
	"sync/atomic"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

// Transaction is trm.Transaction for pgx.Tx.
type Transaction struct {
	tx pgx.Tx

	isActive int64
}

// NewTransaction creates trm.Transaction for pgx.Tx.
func NewTransaction(
	ctx context.Context,
	opts *pgx.TxOptions,
	db *pgxpool.Pool,
) (context.Context, *Transaction, error) {
	tx, err := db.BeginTx(ctx, *opts)
	if err != nil {
		return ctx, nil, err
	}

	tr := &Transaction{tx: tx, isActive: 1}

	go tr.awaitDone(ctx)

	return ctx, tr, nil
}

// NewNestedTransaction creates trm.Transaction for pgx.Tx.
func NewNestedTransaction(
	ctx context.Context,
	tx pgx.Tx,
) (context.Context, *Transaction, error) {
	tx, err := tx.Begin(ctx)
	if err != nil {
		return ctx, nil, err
	}

	tr := &Transaction{tx: tx, isActive: 1}

	go tr.awaitDone(ctx)

	return ctx, tr, nil
}

func (t *Transaction) awaitDone(ctx context.Context) {
	if ctx.Done() == nil {
		return
	}

	<-ctx.Done()

	t.deactivate()
}

// Transaction returns the real transaction pgx.Tx.
// trm.NestedTrFactory returns IsActive as true while trm.Transaction is opened.
func (t *Transaction) Transaction() interface{} {
	return t.tx
}

// Begin nested transaction
func (t *Transaction) Begin(ctx context.Context, _ trm.Settings) (context.Context, trm.Transaction, error) { //nolint:ireturn,nolintlint
	return NewNestedTransaction(ctx, t.tx)
}

// Commit closes the trm.Transaction.
func (t *Transaction) Commit(ctx context.Context) error {
	defer t.deactivate()

	return t.tx.Commit(ctx)
}

// Rollback the trm.Transaction.
func (t *Transaction) Rollback(ctx context.Context) error {
	defer t.deactivate()

	return t.tx.Rollback(ctx)
}

// IsActive returns true if the transaction started but not committed or rolled back.
func (t *Transaction) IsActive() bool {
	return atomic.LoadInt64(&t.isActive) == 1
}

func (t *Transaction) deactivate() {
	atomic.SwapInt64(&t.isActive, 0)
}
