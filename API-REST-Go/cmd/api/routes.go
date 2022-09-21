package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	util "API-REST/cmd/api/utilities"
)

func (s *server) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), struct{ p string }{"params"}, ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *server) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(s.checkToken)

	// Unsecured routes
	router.HandlerFunc(http.MethodGet, "/status", s.statusHandler)

	router.HandlerFunc(http.MethodPost, "/v1/login", s.controllers.user.Login)
	router.HandlerFunc(http.MethodPost, "/v1/register", s.controllers.user.Insert)
	router.HandlerFunc(http.MethodGet, "/v1/users", s.controllers.user.GetAll)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", s.controllers.user.Get)
	router.HandlerFunc(http.MethodPut, "/v1/users/update/:id", s.controllers.user.Update)
	router.HandlerFunc(http.MethodDelete, "/v1/users/delete/:id", s.controllers.user.Delete)

	// ...

	// Secured routes
	router.GET("/v1/assets", s.wrap(secure.ThenFunc(s.controllers.asset.GetAll)))
	// router.POST("/v1/create-payment-intent", s.wrap(secure.ThenFunc(s.createPaymentIntent)))
	// router.POST("/v1/confirm-payment", s.wrap(secure.ThenFunc(s.confirmPayment)))
	// ...

	return s.enableCORS(router)
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

// statusHandler handles /status
func (s *server) statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status:      "Available",
		Environment: s.config.env,
		Version:     VERSION,
	}
	util.WriteJSON(w, http.StatusOK, currentStatus, "status")
}
