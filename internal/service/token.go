package service

import (
	"context"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"github/cntrkilril/auth-golang/internal/controller"
	"github/cntrkilril/auth-golang/internal/entity"
	"github/cntrkilril/auth-golang/pkg/hasher"
	"github/cntrkilril/auth-golang/pkg/tokens"
	"time"
)

type (
	TokenService struct {
		tokenRepo             TokenGateway
		expiresInAccessToken  time.Duration
		expiresInRefreshToken time.Duration
		hasher                hasher.Interactor
		tokensWorker          tokens.Interactor
	}
)

func (s *TokenService) CreateTokens(ctx context.Context, dto controller.CreateTokensDTO) (entity.Tokens, error) {
	accessToken, err := s.tokensWorker.CreateAccessToken(tokens.AccessTokenClaims{
		UserID: dto.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiresInAccessToken)),
		},
	})
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	refreshToken, err := s.tokensWorker.CreateRefreshToken()
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	hashedRefreshToken, err := s.hasher.HashPassword(refreshToken)
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, entity.CreateRefreshTokenParams{
		UserID:    dto.UserID,
		Token:     hashedRefreshToken,
		ExpiresAt: time.Now().Add(s.expiresInRefreshToken),
	})
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	return entity.Tokens{
		AccessToken: entity.AccessToken{
			Token: accessToken,
		},
		RefreshToken: entity.RefreshToken{
			Token: base64.StdEncoding.EncodeToString([]byte(refreshToken)),
		},
	}, nil

}

func (s *TokenService) RefreshToken(ctx context.Context, dto entity.Tokens) (entity.Tokens, error) {

	accessTokenClaims, err := s.tokensWorker.CheckAccessToken(tokens.Token{Token: dto.AccessToken.Token})
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	refreshTokenRes, err := s.tokenRepo.FindRefreshToken(ctx, accessTokenClaims.UserID)
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	refreshTokenByte, err := base64.StdEncoding.DecodeString(dto.RefreshToken.Token)
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	var refreshToken entity.FindRefreshTokenResponse
	var sucFlag bool
	for _, i := range refreshTokenRes {
		if s.hasher.CompareAndHash(i.Token, string(refreshTokenByte)) {
			refreshToken = i
			sucFlag = true
			break
		}
	}

	if !sucFlag {
		return entity.Tokens{}, HandleServiceError(entity.ErrTokenNotFound)
	}

	if accessTokenClaims.UserID != refreshToken.UserID {
		return entity.Tokens{}, HandleServiceError(err)
	}

	err = s.tokensWorker.CheckRefreshToken(
		tokens.RefreshToken{
			UserID:    refreshToken.UserID,
			Token:     refreshToken.Token,
			ExpiresAt: refreshToken.ExpiresAt,
		},
		tokens.Token{Token: dto.RefreshToken.Token},
	)
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	err = s.tokenRepo.DeleteRefreshToken(ctx, refreshToken.Token)
	if err != nil {
		return entity.Tokens{}, err
	}

	accessToken, err := s.tokensWorker.CreateAccessToken(tokens.AccessTokenClaims{
		UserID: refreshToken.UserID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiresInAccessToken)),
		},
	})
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	newRefreshToken, err := s.tokensWorker.CreateRefreshToken()
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	hashedRefreshToken, err := s.hasher.HashPassword(newRefreshToken)
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	err = s.tokenRepo.SaveRefreshToken(ctx, entity.CreateRefreshTokenParams{
		UserID:    refreshToken.UserID,
		Token:     hashedRefreshToken,
		ExpiresAt: time.Now().Add(s.expiresInRefreshToken),
	})
	if err != nil {
		return entity.Tokens{}, HandleServiceError(err)
	}

	return entity.Tokens{
		AccessToken: entity.AccessToken{
			Token: accessToken,
		},
		RefreshToken: entity.RefreshToken{
			Token: newRefreshToken,
		},
	}, nil

}

var _ controller.TokenService = (*TokenService)(nil)

func NewTokenService(
	tokenRepo TokenGateway,
	expiresInAccessToken time.Duration,
	expiresInRefreshToken time.Duration,
	hasher hasher.Interactor,
	tokensWorker tokens.Interactor,
) *TokenService {
	return &TokenService{
		tokenRepo:             tokenRepo,
		expiresInAccessToken:  expiresInAccessToken,
		expiresInRefreshToken: expiresInRefreshToken,
		hasher:                hasher,
		tokensWorker:          tokensWorker,
	}
}
