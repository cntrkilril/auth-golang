package infrastructure

import (
	"context"
	"fmt"
	"github/cntrkilril/auth-golang/internal/entity"
	"github/cntrkilril/auth-golang/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenRepository struct {
	db mongo.Database
}

func (r *TokenRepository) SaveRefreshToken(ctx context.Context, p entity.CreateRefreshTokenParams) error {
	res, err := r.db.Collection("invite_tokens").InsertOne(ctx, bson.M{
		"userID":    p.UserID,
		"token":     p.Token,
		"expiresAt": p.ExpiresAt,
	})
	fmt.Println(res)
	if err != nil {
		return err
	}

	return nil
}

func (r *TokenRepository) FindRefreshToken(ctx context.Context, userID string) (result []entity.FindRefreshTokenResponse, err error) {
	cur, err := r.db.Collection("invite_tokens").Find(ctx, bson.M{"userID": userID})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []entity.FindRefreshTokenResponse{}, entity.ErrTokenNotFound
		}
		return []entity.FindRefreshTokenResponse{}, err
	}
	for cur.Next(ctx) {

		var elem entity.FindRefreshTokenResponse
		err := cur.Decode(&elem)
		if err != nil {
			return []entity.FindRefreshTokenResponse{}, err
		}

		result = append(result, elem)
	}

	return result, err
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Collection("invite_tokens").DeleteOne(ctx, bson.M{
		"token": token,
	})
	if err != nil {
		return err
	}
	return nil
}

var _ service.TokenGateway = (*TokenRepository)(nil)

func NewTokenRepository(db mongo.Database) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}
