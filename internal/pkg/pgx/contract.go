package pgx

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Tr is an interface to work with pgx.Tx.
// Stmtx, StmtxContext, NamedStmt and NamedStmtContext are not implemented!
//
//nolint:interfacebloat
type Tr interface {
	pgx.Tx
}

type Executor interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}
