package utilities

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteJSON(ctx *gin.Context, status int, data interface{}, wrap string) {
	ctx.JSON(status, gin.H{wrap: data})
}

type errorResponse struct {
	Message string `json:"message"`
}

func ErrorJSON(ctx *gin.Context, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	theError := errorResponse{
		Message: err.Error(),
	}

	WriteJSON(ctx, statusCode, theError, "error")
}
