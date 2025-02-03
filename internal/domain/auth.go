package domain

import "context"

type Auth struct {
	ID           int64  `json:"id"`
	ClientID     int64  `json:"client_id"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRepository interface {
	GetAuthByClientId(ctx context.Context, clientId int64) (Auth, error)
	UpdateAuthByClientId(ctx context.Context, token string, refreshToken string, clientId int64) error
}
