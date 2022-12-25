package middleware

import (
	"API-REST/api-gateway/controllers"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"
)

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization") // It tells the client Authorization is important

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			// could set an anonoymous user
			util.ErrorJSON(w, errors.New("invalid authorization header"))
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			util.ErrorJSON(w, errors.New("invalid authorization header"))
			return
		}

		if headerParts[0] != "Bearer" {
			util.ErrorJSON(w, errors.New("unauthorized - no bearer"))
			return
		}

		token := headerParts[1]
		claims, err := jwt.HMACCheck([]byte(token), []byte(conf.JwtSecret))
		if err != nil {
			util.ErrorJSON(w, err, http.StatusUnauthorized)
			return
		}

		if !claims.Valid(time.Now()) {
			util.ErrorJSON(w, errors.New("unauthorized - token expired"), http.StatusUnauthorized)
			return
		}

		if !claims.AcceptAudience(conf.Domain) {
			util.ErrorJSON(w, errors.New("unauthorized - invalid audience"), http.StatusUnauthorized)
			return
		}

		if claims.Issuer != conf.Domain {
			util.ErrorJSON(w, errors.New("unauthorized - invalid issuer"), http.StatusUnauthorized)
			return
		}

		claimerId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			util.ErrorJSON(w, errors.New("unauthorized - invalid claimer"), http.StatusUnauthorized)
			return
		}
		claimerRole, err := controllers.User.CheckRole(claimerId)
		if err != nil {
			util.ErrorJSON(w, errors.New("forbidden - invalid claimer"), http.StatusForbidden)
			return
		}

		// Add claimer ID and claimer Role to header, so we know who is making this request
		ctx := context.WithValue(r.Context(), "Claimer-ID", claims.Subject)
		ctx2 := context.WithValue(ctx, "Claimer-Role", claimerRole)

		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}
