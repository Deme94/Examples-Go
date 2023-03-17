package middleware

import (
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pascaldekloe/jwt"
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
	claims, err := jwt.HMACCheck([]byte(token), []byte(conf.Env.GetString("JWT_AUTH_SECRET")))
	if err != nil {
		return util.ErrorJSON(ctx, errors.New("unauthorized - invalid token"), http.StatusUnauthorized)
	}

	if !claims.Valid(time.Now()) {
		return util.ErrorJSON(ctx, errors.New("unauthorized - token expired"), http.StatusUnauthorized)
	}

	domain := conf.Env.GetString("DOMAIN")
	if !claims.AcceptAudience(domain) {
		return util.ErrorJSON(ctx, errors.New("unauthorized - invalid audience"), http.StatusUnauthorized)
	}

	if claims.Issuer != domain {
		return util.ErrorJSON(ctx, errors.New("unauthorized - invalid issuer"), http.StatusUnauthorized)
	}

	claimerID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return util.ErrorJSON(ctx, errors.New("unauthorized - invalid claimer"), http.StatusUnauthorized)
	}

	// Add claimer ID and claimer roles to header, so we know who is making this request
	ctx.Locals("Claimer-ID", claimerID)

	return ctx.Next()
}
