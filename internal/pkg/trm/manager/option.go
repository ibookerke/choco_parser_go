package manager

import "github.com/ibookerke/choco_parser_go/internal/pkg/trm"

// WithLog sets logger for Manager.
func WithLog(l logger) Opt {
	return func(m *Manager) error {
		m.log = l

		return nil
	}
}

// WithSettings sets trm.Settings for Manager.
func WithSettings(s trm.Settings) Opt {
	return func(m *Manager) error {
		m.settings = s

		return nil
	}
}

// WithCtxManager sets trm.Settings for Manager.
func WithCtxManager(c trm.СtxManager) Opt {
	return func(m *Manager) error {
		m.ctxManager = c

		return nil
	}
}
