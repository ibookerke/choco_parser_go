package repository

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrInternal      = errors.New("internal error")
	ErrNotFound      = errors.New("not found")
)

func wrapScanError(err error) error {
	errStorage := ErrInternal

	if errors.Is(err, pgx.ErrNoRows) {
		errStorage = ErrNotFound
	} else {
		var errPg *pgconn.PgError

		if errors.As(err, &errPg) {
			if errPg.SQLState() == "23505" {
				errStorage = ErrAlreadyExists
			}
		}
	}

	return fmt.Errorf("%w: %s", errStorage, err)
}
