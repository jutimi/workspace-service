package service

import (
	"context"
	"workspace-server/app/model"
)

type OAuthService interface {
	RefreshToken(ctx context.Context, data *model.RefreshTokenRequest) (*model.RefreshTokenResponse, error)
}
