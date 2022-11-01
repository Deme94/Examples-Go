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
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (s *server) routes() http.Handler {
	router := httprouter.New()

	secureUser := alice.New(s.checkToken)
	secureAdmin := alice.New(s.checkToken, s.checkAdmin)

	// Unsecured routes
	router.HandlerFunc(http.MethodGet, "/status", s.statusHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/login", s.controllers.user.Login)
	// swagger:route GET /users users
	//
	// Lists users filtered by some parameters.
	//
	// This will show all available users by default.
	//
	//     Produces:
	//     - application/json
	//
	//     Schemes: http
	//
	//     Parameters:
	//       + name: year
	//         in: query
	//         description: year the user was created
	//         required: false
	//         type: integer
	//         format: int32
	//
	//
	//     Responses:
	//       200: usersResponse
	//     		description: users response
	//     		schema:
	//       		type: array
	//       400: errorResponse
	router.HandlerFunc(http.MethodPost, "/api/v1/users", s.controllers.user.Insert)
	router.HandlerFunc(http.MethodGet, "/api/v1/users", s.controllers.user.GetAll)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:id", s.controllers.user.Get)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:id/photo", s.controllers.user.GetPhoto)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:id/cv", s.controllers.user.GetCV)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/:id", s.controllers.user.Update)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/:id/photo", s.controllers.user.UpdatePhoto)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/:id/cv", s.controllers.user.UpdateCV)
	router.HandlerFunc(http.MethodDelete, "/api/v1/users/:id", s.controllers.user.Delete)

	router.HandlerFunc(http.MethodGet, "/api/v1/assets", s.controllers.asset.GetAll)
	router.HandlerFunc(http.MethodGet, "/api/v1/assets/:id", s.controllers.asset.Get)
	router.HandlerFunc(http.MethodGet, "/api/v1/assets-names", s.controllers.asset.GetNames)
	router.HandlerFunc(http.MethodPost, "/api/v1/assets", s.controllers.asset.Insert)
	router.HandlerFunc(http.MethodPut, "/api/v1/assets/:id", s.controllers.asset.Update)
	router.HandlerFunc(http.MethodDelete, "/api/v1/assets/:id", s.controllers.asset.Delete)

	router.HandlerFunc(http.MethodGet, "/api/v1/attributes", s.controllers.attribute.GetAll)
	router.HandlerFunc(http.MethodGet, "/api/v1/attributes/:id", s.controllers.attribute.Get)
	router.HandlerFunc(http.MethodPost, "/api/v1/attributes", s.controllers.attribute.Insert)
	router.HandlerFunc(http.MethodPut, "/api/v1/attributes/:id", s.controllers.attribute.Update)
	router.HandlerFunc(http.MethodDelete, "/api/v1/attributes/:id", s.controllers.attribute.Delete)

	// ...

	// Secured routes
	router.GET("/api/v1/users/:id/exampleSecuredUserOrAdmin", s.wrap(secureUser.ThenFunc(s.controllers.user.GetSecuredUser)))
	router.GET("/api/v1/users/:id/exampleSecuredOnlyAdmin", s.wrap(secureAdmin.ThenFunc(s.controllers.user.GetSecuredAdmin)))
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
