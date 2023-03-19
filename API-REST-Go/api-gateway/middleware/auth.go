package middleware

import (
	"API-REST/api-gateway/controllers"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CheckToken(ctx *fiber.Ctx) error {
	ctx.Set("Vary", "Authorization") // It tells the client Authorization is important

	authHeader := ctx.Get("Authorization")

	if authHeader == "" {
		// could set an anonoymous user
		return util.ErrorJSON(ctx, errors.New("invalid authorization header"), http.StatusUnauthorized)
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return util.ErrorJSON(ctx, errors.New("invalid authorization header"), http.StatusUnauthorized)
	}

	if headerParts[0] != "Bearer" {
		return util.ErrorJSON(ctx, errors.New("unauthorized - no bearer"), http.StatusUnauthorized)
	}

	token := headerParts[1]
	userID, err := controllers.User.Auth.ValidateJwtToken([]byte(token), conf.Env.GetString("JWT_AUTH_SECRET"))
	if err != nil {
		return util.ErrorJSON(ctx, err, http.StatusUnauthorized)
	}

	// Add claimer ID and claimer roles to header, so we know who is making this request
	ctx.Locals("Claimer-ID", userID.String())

	return ctx.Next()
}
