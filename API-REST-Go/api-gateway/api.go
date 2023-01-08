package api

import (
	"API-REST/api-gateway/controllers"
	"API-REST/api-gateway/middleware"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"API-REST/services/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Start() error {
	// Build controllers
	controllers.Build()

	// Set GIN Mode
	if conf.Env.GetString("ENVIRONMENT") == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup GIN api server
	r := gin.Default()
	r.MaxMultipartMemory = conf.Conf.GetInt64("maxMultiPartMemory")

	// Setup groups and middleware
	api := r.Group(conf.Conf.GetString("apiBasePath"))
	api.Use(middleware.EnableCORS)

	apiPrivate := api.Group("/")
	apiPrivate.Use(middleware.CheckToken)

	apiAdmin := apiPrivate.Group("/")
	apiAdmin.Use(middleware.CheckAdmin)

	api.GET("/status", func(ctx *gin.Context) {
		appStatus := struct {
			Status      string `json:"status"`
			Environment string `json:"environment"`
			Version     string `json:"version"`
		}{
			Status:      "Available",
			Environment: conf.Env.GetString("ENVIRONMENT"),
			Version:     conf.Env.GetString("VERSION"),
		}

		util.WriteJSON(ctx, http.StatusOK, appStatus, "status")
	})

	// Unsecured routes
	api.POST("/login", controllers.User.Login)
	api.POST("/users", controllers.User.Insert)
	api.GET("/users", controllers.User.GetAll)
	api.GET("/users/:id", controllers.User.Get)
	api.GET("/users/:id/photo", controllers.User.GetPhoto)
	api.GET("/users/:id/cv", controllers.User.GetCV)
	api.PUT("/users/:id", controllers.User.Update)
	api.POST("/users/:id/roles", controllers.User.UpdateRoles)
	api.PUT("/users/:id/photo", controllers.User.UpdatePhoto)
	api.PUT("/users/:id/cv", controllers.User.UpdateCV)
	api.DELETE("/users/:id", controllers.User.Delete)

	api.GET("/roles", controllers.Role.GetAll)
	api.GET("/roles/:id", controllers.Role.Get)
	api.POST("/roles", controllers.Role.Insert)
	api.POST("/roles/:id", controllers.Role.Update)
	api.POST("/roles/:id/permissions", controllers.Role.UpdatePermissions)
	api.DELETE("/roles/:id", controllers.Role.Delete)

	api.GET("/permissions", controllers.Permission.GetAll)
	api.GET("/permissions/:id", controllers.Permission.Get)
	api.POST("/permissions", controllers.Permission.Insert)
	api.PUT("/permissions/:id", controllers.Permission.Update)
	api.DELETE("/permissions/:id", controllers.Permission.Delete)

	api.GET("/assets", controllers.Asset.GetAll)
	api.GET("/assets/:id", controllers.Asset.Get)
	api.GET("/assets/:id/attributes", controllers.Asset.GetWithAttributes)
	api.GET("/assets-names", controllers.Asset.GetNames)
	api.POST("/assets", controllers.Asset.Insert)
	api.PUT("/assets/:id", controllers.Asset.Update)
	api.DELETE("/assets/:id", controllers.Asset.Delete)

	api.GET("/attributes", controllers.Attribute.GetAll)
	api.GET("/attributes/:id", controllers.Attribute.Get)
	api.POST("/attributes", controllers.Attribute.Insert)
	api.PUT("/attributes/:id", controllers.Attribute.Update)
	api.DELETE("/attributes/:id", controllers.Attribute.Delete)
	// ...

	// Secured routes
	apiPrivate.GET("/users/:id/exampleSecuredUserOrAdmin", controllers.User.GetSecuredUser)
	apiAdmin.GET("/users/:id/exampleSecuredOnlyAdmin", controllers.User.GetSecuredAdmin)
	// ...

	// Log start
	host := conf.Env.GetString("HOST")
	port := conf.Env.GetString("PORT")
	logger.Logger.Printf("Starting server on http://%s:%s\n", host, port)

	// Run server
	return r.Run(conf.Env.GetString("HOST") + ":" + port)
}
