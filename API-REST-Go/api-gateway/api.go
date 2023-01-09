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

	apiAuth := api.Group("/")
	apiAuth.Use(middleware.CheckToken)

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

	api.POST("/auth", controllers.Auth.Login)
	apiAuth.GET("/auth", controllers.Auth.Get)
	apiAuth.GET("/auth/photo", controllers.Auth.GetPhoto)
	apiAuth.GET("/auth/cv", controllers.Auth.GetCV)
	apiAuth.PUT("/auth", controllers.Auth.Update)
	apiAuth.PUT("/auth/photo", controllers.Auth.UpdatePhoto)
	apiAuth.PUT("/auth/cv", controllers.Auth.UpdateCV)
	apiAuth.DELETE("/auth", controllers.Auth.Delete)

	apiAuth.GET("/users", middleware.CheckPermission("users", "read"), controllers.User.GetAll)
	apiAuth.GET("/users/:id", middleware.CheckPermission("users", "read"), controllers.User.Get)
	apiAuth.GET("/users/:id/photo", middleware.CheckPermission("users", "read"), controllers.User.GetPhoto)
	apiAuth.GET("/users/:id/cv", middleware.CheckPermission("users", "read"), controllers.User.GetCV)
	api.POST("/users", controllers.User.Insert) // public registration
	apiAuth.POST("/users/:id/roles", middleware.CheckPermission("users", "assign"), controllers.User.UpdateRoles)
	apiAuth.PUT("/users/:id", middleware.CheckPermission("users", "update"), controllers.User.Update)
	apiAuth.PUT("/users/:id/photo", middleware.CheckPermission("users", "update"), controllers.User.UpdatePhoto)
	apiAuth.PUT("/users/:id/cv", middleware.CheckPermission("users", "update"), controllers.User.UpdateCV)
	apiAuth.DELETE("/users/:id", middleware.CheckPermission("users", "delete"), controllers.User.Delete)

	apiAuth.GET("/roles", middleware.CheckPermission("roles", "read"), controllers.Role.GetAll)
	apiAuth.GET("/roles/:id", middleware.CheckPermission("roles", "read"), controllers.Role.Get)
	apiAuth.POST("/roles", middleware.CheckPermission("roles", "create"), controllers.Role.Insert)
	apiAuth.POST("/roles/:id/permissions", middleware.CheckPermission("roles", "update"), controllers.Role.UpdatePermissions)
	apiAuth.PUT("/roles/:id", middleware.CheckPermission("roles", "update"), controllers.Role.Update)
	apiAuth.DELETE("/roles/:id", middleware.CheckPermission("roles", "delete"), controllers.Role.Delete)

	apiAuth.GET("/permissions", middleware.CheckPermission("permissions", "read"), controllers.Permission.GetAll)
	apiAuth.GET("/permissions/:id", middleware.CheckPermission("permissions", "read"), controllers.Permission.Get)
	apiAuth.POST("/permissions", middleware.CheckPermission("permissions", "create"), controllers.Permission.Insert)
	apiAuth.PUT("/permissions/:id", middleware.CheckPermission("permissions", "update"), controllers.Permission.Update)
	apiAuth.DELETE("/permissions/:id", middleware.CheckPermission("permissions", "delete"), controllers.Permission.Delete)

	apiAuth.GET("/assets", middleware.CheckPermission("assets", "read"), controllers.Asset.GetAll)
	apiAuth.GET("/assets/:id", middleware.CheckPermission("assets", "read"), controllers.Asset.Get)
	apiAuth.GET("/assets/:id/attributes", middleware.CheckPermission("assets", "read"), controllers.Asset.GetWithAttributes)
	apiAuth.GET("/assets/names", middleware.CheckPermission("assets", "read"), controllers.Asset.GetNames)
	apiAuth.POST("/assets", middleware.CheckPermission("assets", "create"), controllers.Asset.Insert)
	apiAuth.PUT("/assets/:id", middleware.CheckPermission("assets", "update"), controllers.Asset.Update)
	apiAuth.DELETE("/assets/:id", middleware.CheckPermission("assets", "delete"), controllers.Asset.Delete)

	apiAuth.GET("/attributes", middleware.CheckPermission("attributes", "read"), controllers.Attribute.GetAll)
	apiAuth.GET("/attributes/:id", middleware.CheckPermission("attributes", "read"), controllers.Attribute.Get)
	apiAuth.POST("/attributes", middleware.CheckPermission("attributes", "create"), controllers.Attribute.Insert)
	apiAuth.PUT("/attributes/:id", middleware.CheckPermission("attributes", "update"), controllers.Attribute.Update)
	apiAuth.DELETE("/attributes/:id", middleware.CheckPermission("attributes", "delete"), controllers.Attribute.Delete)
	// ...

	// Log start
	host := conf.Env.GetString("HOST")
	port := conf.Env.GetString("PORT")
	logger.Logger.Printf("Starting server on http://%s:%s\n", host, port)

	// Run server
	return r.Run(conf.Env.GetString("HOST") + ":" + port)
}
