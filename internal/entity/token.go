package entity

import "time"

type (
	AccessToken struct {
		Token string `json:"accessToken" validate:"required"`
	}

	RefreshToken struct {
		Token string `json:"refreshToken" validate:"required"`
	}

	Tokens struct {
		AccessToken
		RefreshToken
	}

	CreateRefreshTokenParams struct {
		UserID    string
		Token     string
		ExpiresAt time.Time
	}

	FindRefreshTokenResponse struct {
		UserID    string
		Token     string
		ExpiresAt time.Time
	}
)
