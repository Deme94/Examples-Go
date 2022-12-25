package api

import (
	"API-REST/api-gateway/controllers"
	"API-REST/api-gateway/middleware"
	util "API-REST/api-gateway/utilities"
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func Routes() http.Handler {
	router := httprouter.New()

	secureUser := alice.New(middleware.CheckToken)
	secureAdmin := alice.New(middleware.CheckToken, middleware.CheckAdmin)

	// Unsecured routes
	router.HandlerFunc(http.MethodGet, "/status", statusHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/login", controllers.User.Login)
	router.HandlerFunc(http.MethodPost, "/api/v1/users", controllers.User.Insert)
	router.HandlerFunc(http.MethodGet, "/api/v1/users", controllers.User.GetAll)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:id", controllers.User.Get)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:id/photo", controllers.User.GetPhoto)
	router.HandlerFunc(http.MethodGet, "/api/v1/users/:id/cv", controllers.User.GetCV)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/:id", controllers.User.Update)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/:id/photo", controllers.User.UpdatePhoto)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/:id/cv", controllers.User.UpdateCV)
	router.HandlerFunc(http.MethodDelete, "/api/v1/users/:id", controllers.User.Delete)

	router.HandlerFunc(http.MethodGet, "/api/v1/assets", controllers.Asset.GetAll)
	router.HandlerFunc(http.MethodGet, "/api/v1/assets/:id", controllers.Asset.Get)
	router.HandlerFunc(http.MethodGet, "/api/v1/assets/:id/attributes", controllers.Asset.GetWithAttributes)
	router.HandlerFunc(http.MethodGet, "/api/v1/assets-names", controllers.Asset.GetNames)
	router.HandlerFunc(http.MethodPost, "/api/v1/assets", controllers.Asset.Insert)
	router.HandlerFunc(http.MethodPut, "/api/v1/assets/:id", controllers.Asset.Update)
	router.HandlerFunc(http.MethodDelete, "/api/v1/assets/:id", controllers.Asset.Delete)

	router.HandlerFunc(http.MethodGet, "/api/v1/attributes", controllers.Attribute.GetAll)
	router.HandlerFunc(http.MethodGet, "/api/v1/attributes/:id", controllers.Attribute.Get)
	router.HandlerFunc(http.MethodPost, "/api/v1/attributes", controllers.Attribute.Insert)
	router.HandlerFunc(http.MethodPut, "/api/v1/attributes/:id", controllers.Attribute.Update)
	router.HandlerFunc(http.MethodDelete, "/api/v1/attributes/:id", controllers.Attribute.Delete)

	// ...

	// Secured routes
	router.GET("/api/v1/users/:id/exampleSecuredUserOrAdmin", wrap(secureUser.ThenFunc(controllers.User.GetSecuredUser)))
	router.GET("/api/v1/users/:id/exampleSecuredOnlyAdmin", wrap(secureAdmin.ThenFunc(controllers.User.GetSecuredAdmin)))
	// router.POST("/v1/create-payment-intent", wrap(secure.ThenFunc(createPaymentIntent)))
	// router.POST("/v1/confirm-payment", wrap(secure.ThenFunc(confirmPayment)))
	// ...

	return middleware.EnableCORS(router)
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

// statusHandler handles /status
func statusHandler(w http.ResponseWriter, r *http.Request) {
	currentStatus := AppStatus{
		Status: "Available",
		//Environment: S.Config.env,
		//Version:     VERSION,
	}
	util.WriteJSON(w, http.StatusOK, currentStatus, "status")
}
