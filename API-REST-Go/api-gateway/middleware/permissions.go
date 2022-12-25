package middleware

import (
	util "API-REST/api-gateway/utilities"
	"API-REST/services/logger"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckAdmin(ctx *gin.Context) {

	claimerId := ctx.GetString("Claimer-ID")
	claimerRole := ctx.GetString("Claimer-Role")

	if claimerRole != "admin" {
		logger.Logger.Println("USER TRIED TO ACCESS ADMIN OPERATION id:", claimerId)
		util.ErrorJSON(ctx, errors.New("unauthorized - admin required"), http.StatusForbidden)
		return
	}

	ctx.Next()
}
