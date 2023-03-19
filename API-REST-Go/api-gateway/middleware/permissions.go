package middleware

import (
	"API-REST/api-gateway/controllers"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/logger"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CheckPermission(resource string, operation string) func(ctx *fiber.Ctx) error {

	return func(ctx *fiber.Ctx) error {
		claimerID := ctx.Locals("Claimer-ID").(string)
		claimerRoles, err := controllers.User.Auth.GetRoles(uuid.MustParse(claimerID))
		if err != nil {
			return util.ErrorJSON(ctx, err, http.StatusForbidden)
		}

		for _, role := range claimerRoles {
			if role == "superadmin" {
				return ctx.Next()
			}
		}

		hasPerm, err := controllers.User.Auth.HasPermission(uuid.MustParse(claimerID), resource, operation)
		if err != nil {
			return util.ErrorJSON(ctx, err, http.StatusForbidden)
		}

		if !hasPerm {
			logger.Logger.Println("USER TRIED TO ACCESS RESTRINGED OPERATION id:", claimerID)
			return util.ErrorJSON(ctx, errors.New("forbidden operation - user has not permission"), http.StatusForbidden)
		}

		return ctx.Next()
	}
}
