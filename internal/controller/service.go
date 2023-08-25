package controller

import (
	"context"
	"github/cntrkilril/auth-golang/internal/entity"
)

type (
	TokenService interface {
		CreateTokens(context.Context, CreateTokensDTO) (entity.Tokens, error)
		RefreshToken(context.Context, entity.Tokens) (entity.Tokens, error)
	}

	CreateTokensDTO struct {
		UserID string `json:"userID" validate:"required"`
	}
)
