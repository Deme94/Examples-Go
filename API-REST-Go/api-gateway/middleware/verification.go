package middleware

import (
	"API-REST/api-gateway/controllers"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func VerifyEmail(ctx *fiber.Ctx) error {
	if conf.Conf.GetBool("verifyEmail") {
		claimerID := ctx.Locals("Claimer-ID").(int)
		verifiedEmail, err := controllers.User.Auth.HasVerifiedEmail(claimerID)
		if err != nil {
			return util.ErrorJSON(ctx, err, http.StatusInternalServerError)
		}

		if !verifiedEmail {
			return util.ErrorJSON(ctx, errors.New("email is not verified"), http.StatusUnauthorized)
		}
	}
	return ctx.Next()
}
