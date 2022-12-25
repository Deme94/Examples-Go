package middleware

import (
	util "API-REST/api-gateway/utilities"
	"API-REST/services/logger"
	"errors"
	"net/http"
)

func CheckAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claimerId := r.Context().Value("Claimer-ID").(string)
		claimerRole := r.Context().Value("Claimer-Role").(string)

		if claimerRole != "admin" {
			logger.Logger.Println("USER TRIED TO ACCESS ADMIN OPERATION id:", claimerId)
			util.ErrorJSON(w, errors.New("unauthorized - admin required"), http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
