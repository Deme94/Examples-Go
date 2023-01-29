package utilities

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func WriteJSON(ctx *fiber.Ctx, status int, data interface{}, wrap string) error {
	return ctx.Status(status).JSON(map[string]interface{}{wrap: data})
}

type errorResponse struct {
	Message string `json:"message"`
}

func ErrorJSON(ctx *fiber.Ctx, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	theError := errorResponse{
		Message: err.Error(),
	}

	return WriteJSON(ctx, statusCode, theError, "error")
}
