package middleware

import (
	"API-REST/api-gateway/controllers"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/logger"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckPermission(resource string, operation string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		claimerID := ctx.GetInt("Claimer-ID")
		claimerRoles, err := controllers.Auth.GetRoles(claimerID)
		if err != nil {
			util.ErrorJSON(ctx, err, http.StatusForbidden)
			ctx.Abort()
			return
		}

		for _, role := range claimerRoles {
			if role == "admin" {
				ctx.Next()
				return
			}
		}

		hasPerm, err := controllers.Auth.HasPermission(claimerID, resource, operation)
		if err != nil {
			util.ErrorJSON(ctx, err, http.StatusForbidden)
			ctx.Abort()
			return
		}

		if !hasPerm {
			logger.Logger.Println("USER TRIED TO ACCESS RESTRINGED OPERATION id:", claimerID)
			util.ErrorJSON(ctx, errors.New("forbidden operation - user has not permission"), http.StatusForbidden)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
