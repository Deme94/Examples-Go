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
	claimerRoles := ctx.GetStringSlice("Claimer-Roles")

	for _, role := range claimerRoles {
		if role == "admin" {
			ctx.Next()
			return
		}
	}

	logger.Logger.Println("USER TRIED TO ACCESS ADMIN OPERATION id:", claimerId)
	util.ErrorJSON(ctx, errors.New("unauthorized - admin required"), http.StatusForbidden)
}
