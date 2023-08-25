package v1

import "github.com/gofiber/fiber/v2"

func newResp(data, metadata any) fiber.Map {
	return fiber.Map{
		"metadata": metadata,
		"data":     data,
	}
}

func newErrResp(err error) fiber.Map {
	return fiber.Map{
		"error": fiber.Map{
			"message": err.Error(),
		},
	}
}
