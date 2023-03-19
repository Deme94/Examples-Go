package middleware

import (
	"API-REST/api-gateway/controllers"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func VerifyEmail(ctx *fiber.Ctx) error {
	if conf.Conf.GetBool("verifyEmail") {
		claimerID := ctx.Locals("Claimer-ID").(string)
		verifiedEmail, err := controllers.User.Auth.HasVerifiedEmail(uuid.MustParse(claimerID))
		if err != nil {
			return util.ErrorJSON(ctx, err, http.StatusInternalServerError)
		}

		if !verifiedEmail {
			return util.ErrorJSON(ctx, errors.New("email is not verified"), http.StatusUnauthorized)
		}
	}
	return ctx.Next()
}
