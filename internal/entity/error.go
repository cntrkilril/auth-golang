package entity

import "errors"

type ErrCode int

const (
	_ = ErrCode(iota)
	ErrCodeBadRequest
	ErrCodeNotFound
)

type Error struct {
	msg  string
	code ErrCode
}

func (e *Error) Error() string {
	return e.msg
}

func (e *Error) Code() ErrCode {
	return e.code
}

func NewError(msg string, code ErrCode) *Error {
	return &Error{msg, code}
}

var (
	ErrTokenNotFound   = errors.New("токен не найден")
	ErrValidationError = errors.New("ошибка валидации")
	ErrUnknown         = errors.New("неизвестная ошибка")
)
