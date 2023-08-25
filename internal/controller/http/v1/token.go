package v1

import (
	"github.com/gofiber/fiber/v2"
	"github/cntrkilril/auth-golang/internal/controller"
	"github/cntrkilril/auth-golang/internal/entity"
	"github/cntrkilril/auth-golang/pkg/govalidator"
)

type TokenHandler struct {
	tokenService controller.TokenService
	val          *govalidator.Validator
}

func (h *TokenHandler) create() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p controller.CreateTokensDTO
		if err := h.val.ValidateParams(c, &p); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		res, err := h.tokenService.CreateTokens(c.Context(), p)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"accessToken":  res.AccessToken,
				"refreshToken": res.RefreshToken,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *TokenHandler) refresh() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var p entity.Tokens
		if err := h.val.ValidateRequestBody(c, &p); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		res, err := h.tokenService.RefreshToken(c.Context(), p)
		if err != nil {
			return HandleError(c, err)
		}
		return c.Status(fiber.StatusOK).JSON(newResp(
			fiber.Map{
				"accessToken":  res.AccessToken,
				"refreshToken": res.RefreshToken,
			},
			fiber.Map{
				"success": true,
			},
		))
	}
}

func (h *TokenHandler) Register(r fiber.Router) {
	r.Get("create-tokens/:userID",
		h.create())
	r.Post("refresh-tokens",
		h.refresh())
}

func NewTokenHandler(
	tokenService controller.TokenService,
	val *govalidator.Validator,
) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
		val:          val,
	}
}
