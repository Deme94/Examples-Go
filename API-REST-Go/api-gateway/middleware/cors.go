package middleware

import (
	"github.com/gin-gonic/gin"
)

func EnableCORS(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

	ctx.Next()
}
