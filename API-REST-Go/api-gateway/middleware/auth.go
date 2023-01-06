package middleware

import (
	"API-REST/api-gateway/controllers"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pascaldekloe/jwt"
)

func CheckToken(ctx *gin.Context) {
	ctx.Header("Vary", "Authorization") // It tells the client Authorization is important

	authHeader := ctx.Request.Header.Get("Authorization")

	if authHeader == "" {
		// could set an anonoymous user
		util.ErrorJSON(ctx, errors.New("invalid authorization header"))
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		util.ErrorJSON(ctx, errors.New("invalid authorization header"))
		return
	}

	if headerParts[0] != "Bearer" {
		util.ErrorJSON(ctx, errors.New("unauthorized - no bearer"))
		return
	}

	token := headerParts[1]
	claims, err := jwt.HMACCheck([]byte(token), []byte(conf.Env.GetString("JWT_SECRET")))
	if err != nil {
		util.ErrorJSON(ctx, err, http.StatusUnauthorized)
		return
	}

	if !claims.Valid(time.Now()) {
		util.ErrorJSON(ctx, errors.New("unauthorized - token expired"), http.StatusUnauthorized)
		return
	}

	domain := conf.Env.GetString("DOMAIN")
	if !claims.AcceptAudience(domain) {
		util.ErrorJSON(ctx, errors.New("unauthorized - invalid audience"), http.StatusUnauthorized)
		return
	}

	if claims.Issuer != domain {
		util.ErrorJSON(ctx, errors.New("unauthorized - invalid issuer"), http.StatusUnauthorized)
		return
	}

	claimerId, err := strconv.Atoi(claims.Subject)
	if err != nil {
		util.ErrorJSON(ctx, errors.New("unauthorized - invalid claimer"), http.StatusUnauthorized)
		return
	}
	claimerRoles, err := controllers.User.CheckRoles(claimerId)
	if err != nil {
		util.ErrorJSON(ctx, errors.New("forbidden - invalid claimer"), http.StatusForbidden)
		return
	}

	// Add claimer ID and claimer roles to header, so we know who is making this request
	ctx.Set("Claimer-ID", claims.Subject)
	ctx.Set("Claimer-Roles", claimerRoles)

	ctx.Next()
}
