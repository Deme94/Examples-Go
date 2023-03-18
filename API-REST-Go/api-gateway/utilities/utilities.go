package utilities

import (
	"net/http"

	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
)

func WriteJSON(ctx *fiber.Ctx, status int, data interface{}, wrap ...string) error {
	if len(wrap) > 0 {
		data = map[string]interface{}{wrap[0]: data}
	}
	return ctx.Status(status).JSON(data)
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

func StructToMap(object interface{}) map[string]interface{} {
	var res map[string]interface{}
	objectBytes, _ := json.Marshal(object)
	json.Unmarshal(objectBytes, &res)

	return res
}
