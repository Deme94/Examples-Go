package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"
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

		claimerId, err := strconv.Atoi(claims.Subject)
		if err != nil {
			util.ErrorJSON(w, errors.New("unauthorized - invalid claimer"), http.StatusForbidden)
			return
		}
		claimerRole, err := s.controllers.user.Model.GetRole(claimerId)
		if err != nil {
			util.ErrorJSON(w, errors.New("unauthorized - invalid claimer"), http.StatusForbidden)
			return
		}

		// Add claimer ID and claimer Role to header, so we know who is making this request
		ctx := context.WithValue(r.Context(), "Claimer-ID", claims.Subject)
		ctx2 := context.WithValue(ctx, "Claimer-Role", claimerRole)

		s.logger.Println("Valid user:", claims.Subject)

		next.ServeHTTP(w, r.WithContext(ctx2))
	})
}
func (s *server) checkAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claimerId := r.Context().Value("Claimer-ID").(string)
		claimerRole := r.Context().Value("Claimer-Role").(string)

		if claimerRole != "admin" {
			s.logger.Println("USER TRIED TO ACCESS ADMIN OPERATION id:", claimerId)
			util.ErrorJSON(w, errors.New("unauthorized - admin required"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
