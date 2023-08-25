package service

import "github/cntrkilril/auth-golang/internal/entity"

func HandleServiceError(err error) error {
	switch err {
	case entity.ErrTokenNotFound:
		return entity.NewError(entity.ErrTokenNotFound.Error(), entity.ErrCodeNotFound)
	case entity.ErrValidationError:
		return entity.NewError(entity.ErrValidationError.Error(), entity.ErrCodeBadRequest)
	default:
		return entity.NewError(entity.ErrUnknown.Error(), entity.ErrCodeBadRequest)
	}
}
