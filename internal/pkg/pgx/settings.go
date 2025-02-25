package pgx

import (
	"github.com/jackc/pgx/v5"

	"github.com/ibookerke/choco_parser_go/internal/pkg/trm"
)

// Opt is a type to configure Settings.
type Opt func(*Settings) error

// WithTxOptions sets up sql.TxOptions for the Settings.
func WithTxOptions(opts *pgx.TxOptions) Opt {
	return func(s *Settings) error {
		*s = s.setTrOpts(opts)

		return nil
	}
}

// Settings contains settings for sql.Transaction.
type Settings struct {
	trm.Settings
	txOpts *pgx.TxOptions
}

// NewSettings creates Settings.
func NewSettings(trms trm.Settings, oo ...Opt) (Settings, error) {
	s := &Settings{Settings: trms, txOpts: nil}

	for _, o := range oo {
		if err := o(s); err != nil {
			return Settings{}, err
		}
	}

	return *s, nil
}

// MustSettings returns Settings if err is nil and panics otherwise.
func MustSettings(trms trm.Settings, oo ...Opt) Settings {
	s, err := NewSettings(trms, oo...)
	if err != nil {
		panic(err)
	}

	return s
}

//revive:disable:exported
func (s Settings) EnrichBy(in trm.Settings) trm.Settings { //nolint:ireturn,nolintlint
	external, ok := in.(Settings)
	if ok {
		if s.TxOpts() == nil {
			s = s.setTrOpts(external.TxOpts())
		}
	}

	s.Settings = s.Settings.EnrichBy(in)

	return s
}

// TxOpts returns trm.CtxKey for the trm.Transaction.
func (s Settings) TxOpts() *pgx.TxOptions {
	if s.txOpts == nil {
		return &pgx.TxOptions{
			IsoLevel:   pgx.ReadCommitted,
			AccessMode: pgx.ReadWrite,
		}
	}
	return s.txOpts
}

func (s Settings) setTrOpts(opts *pgx.TxOptions) Settings {
	s.txOpts = opts

	return s
}
