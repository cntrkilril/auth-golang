package service

import (
	"context"
	"github/cntrkilril/auth-golang/internal/entity"
)

type (
	TokenGateway interface {
		SaveRefreshToken(context.Context, entity.CreateRefreshTokenParams) error
		FindRefreshToken(context.Context, string) ([]entity.FindRefreshTokenResponse, error)
		DeleteRefreshToken(context.Context, string) error
	}
)
