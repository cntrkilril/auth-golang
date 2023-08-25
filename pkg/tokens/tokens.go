package tokens

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github/cntrkilril/auth-golang/pkg/hasher"
	"time"
)

type (
	AccessTokenClaims struct {
		UserID string
		jwt.RegisteredClaims
	}

	Tokener struct {
		jwtKey []byte
		hasher hasher.Interactor
	}

	Token struct {
		Token string
	}

	RefreshToken struct {
		UserID    string
		Token     string
		ExpiresAt time.Time
	}

	Interactor interface {
		CreateAccessToken(claims AccessTokenClaims) (string, error)
		CreateRefreshToken() (string, error)
		CheckAccessToken(Token) (AccessTokenClaims, error)
		CheckRefreshToken(RefreshToken, Token) error
	}
)

func (t *Tokener) CheckAccessToken(token Token) (res AccessTokenClaims, err error) {
	accessTokenJWT, err := jwt.ParseWithClaims(token.Token, &res, func(token *jwt.Token) (interface{}, error) {
		return t.jwtKey, nil
	})

	if err != nil || !accessTokenJWT.Valid {
		return AccessTokenClaims{}, err
	}

	return res, nil
}

func (t *Tokener) CheckRefreshToken(hashedToken RefreshToken, encodedToken Token) error {

	decodedRefreshToken, err := base64.StdEncoding.DecodeString(encodedToken.Token)
	if err != nil {
		return err
	}

	if !t.hasher.CompareAndHash(hashedToken.Token, string(decodedRefreshToken)) {
		return errors.New("bad refresh token")
	}

	loc, _ := time.LoadLocation("Etc/UTC")
	fmt.Println(time.Now().In(loc))
	if hashedToken.ExpiresAt.Before(time.Now().In(loc)) {
		return errors.New("refresh token expires")
	}
	return nil
}

func (t *Tokener) CreateAccessToken(claims AccessTokenClaims) (string, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(t.jwtKey)
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func (t *Tokener) CreateRefreshToken() (string, error) {

	refreshTokenUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return refreshTokenUUID.String(), nil
}

var _ Interactor = (*Tokener)(nil)

func New(
	jwtKey string,
	hasher hasher.Interactor,
) *Tokener {
	return &Tokener{
		jwtKey: []byte(jwtKey),
		hasher: hasher,
	}
}
