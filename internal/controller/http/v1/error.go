package v1

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github/cntrkilril/auth-golang/internal/entity"
)

var errCodeMap = map[entity.ErrCode]int{
	entity.ErrCodeBadRequest: fiber.StatusBadRequest,
	entity.ErrCodeNotFound:   fiber.StatusNotFound,
}

func HandleError(ctx *fiber.Ctx, err error) error {
	appErr := &entity.Error{}
	if errors.As(err, &appErr) {
		c, ok := errCodeMap[appErr.Code()]
		if !ok {
			c = fiber.StatusInternalServerError
		}

		return ctx.Status(c).JSON(newErrResp(appErr))
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(newErrResp(err))
}
