package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/pascaldekloe/jwt"

	util "API-REST/cmd/api/utilities"
)

func (s *server) enableCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		handler.ServeHTTP(w, r)
	})
}

func (s *server) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

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
		claims, err := jwt.HMACCheck([]byte(token), []byte(s.config.jwt.secret))
		if err != nil {
			util.ErrorJSON(w, err, http.StatusForbidden)
			return
		}

		if !claims.Valid(time.Now()) {
			util.ErrorJSON(w, errors.New("unauthorized - token expired"), http.StatusForbidden)
			return
		}

		if !claims.AcceptAudience(domain) {
			util.ErrorJSON(w, errors.New("unauthorized - invalid audience"), http.StatusForbidden)
			return
		}

		if claims.Issuer != domain {
			util.ErrorJSON(w, errors.New("unauthorized - invalid issuer"), http.StatusForbidden)
			return
		}

		s.logger.Println("Valid user:", claims.Subject)

		next.ServeHTTP(w, r)
	})
}
