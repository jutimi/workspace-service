package service

import (
	"context"
	"oauth-server/app/model"
)

type OAuthService interface {
	RefreshToken(ctx context.Context, data *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error)
}
