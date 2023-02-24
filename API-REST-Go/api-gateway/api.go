package api

import (
	"API-REST/api-gateway/controllers"
	"API-REST/api-gateway/middleware"
	util "API-REST/api-gateway/utilities"
	"API-REST/services/conf"
	"API-REST/services/logger"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func Start() error {
	// Build controllers
	controllers.Build()

	// Setup Fiber api server
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Setup groups and middleware
	api := app.Group(conf.Conf.GetString("apiBasePath"))
	api.Use(middleware.CORS())

	pub := api.Group("/public")
	pvt := api.Group("/private")
	pvt.Use(middleware.CheckToken)
	admin := pvt.Group("/admin")
	admin.Use(
		middleware.CheckPermission("*", "*"),
		// Serve all files of project
		filesystem.New(filesystem.Config{
			Root:   http.Dir("./"),
			Browse: true,
		}),
	)

	api.Get("/status", func(ctx *fiber.Ctx) error {
		appStatus := struct {
			Status      string `json:"status"`
			Environment string `json:"environment"`
			Version     string `json:"version"`
		}{
			Status:      "Available",
			Environment: conf.Env.GetString("ENVIRONMENT"),
			Version:     conf.Env.GetString("VERSION"),
		}

		return util.WriteJSON(ctx, http.StatusOK, appStatus, "status")
	})

	pub.Post("/auth/login", controllers.User.Auth.Login)
	pvt.Get("/auth", controllers.User.Auth.Get)
	pvt.Get("/auth/photo", controllers.User.Auth.GetPhoto)
	pvt.Get("/auth/cv", controllers.User.Auth.GetCV)
	pvt.Put("/auth", controllers.User.Auth.Update)
	pvt.Put("/auth/change-password", controllers.User.Auth.ChangePassword)
	pub.Put("/auth/reset-password", controllers.User.Auth.ResetPassword)
	pvt.Put("/auth/photo", controllers.User.Auth.UpdatePhoto)
	pvt.Put("/auth/cv", controllers.User.Auth.UpdateCV)
	pvt.Delete("/auth", controllers.User.Auth.Delete)

	pvt.Get("/users", middleware.CheckPermission("users", "read"), controllers.User.GetAll)
	pvt.Get("/users/:id", middleware.CheckPermission("users", "read"), controllers.User.Get)
	pvt.Get("/users/:id/photo", middleware.CheckPermission("users", "read"), controllers.User.GetPhoto)
	pvt.Get("/users/:id/cv", middleware.CheckPermission("users", "read"), controllers.User.GetCV)
	pub.Post("/users", controllers.User.Insert) // public registration
	pvt.Post("/users/:id/roles", middleware.CheckPermission("users", "assign"), controllers.User.UpdateRoles)
	pvt.Put("/users/:id", middleware.CheckPermission("users", "update"), controllers.User.Update)
	pvt.Put("/users/:id/photo", middleware.CheckPermission("users", "update"), controllers.User.UpdatePhoto)
	pvt.Put("/users/:id/cv", middleware.CheckPermission("users", "update"), controllers.User.UpdateCV)
	pvt.Put("/users/:id/ban", middleware.CheckPermission("users", "ban"), controllers.User.Ban)
	pvt.Put("/users/:id/unban", middleware.CheckPermission("users", "ban"), controllers.User.Unban)
	pvt.Put("/users/:id/restore", middleware.CheckPermission("users", "delete"), controllers.User.Restore)
	pvt.Delete("/users/:id", middleware.CheckPermission("users", "delete"), controllers.User.Delete)

	pvt.Get("/roles", middleware.CheckPermission("roles", "read"), controllers.Role.GetAll)
	pvt.Get("/roles/:id", middleware.CheckPermission("roles", "read"), controllers.Role.Get)
	pvt.Post("/roles", middleware.CheckPermission("roles", "create"), controllers.Role.Insert)
	pvt.Post("/roles/:id/permissions", middleware.CheckPermission("roles", "assign"), controllers.Role.UpdatePermissions)
	pvt.Put("/roles/:id", middleware.CheckPermission("roles", "update"), controllers.Role.Update)
	pvt.Delete("/roles/:id", middleware.CheckPermission("roles", "delete"), controllers.Role.Delete)

	pvt.Get("/permissions", middleware.CheckPermission("permissions", "read"), controllers.Permission.GetAll)
	pvt.Get("/permissions/:id", middleware.CheckPermission("permissions", "read"), controllers.Permission.Get)
	pvt.Post("/permissions", middleware.CheckPermission("permissions", "create"), controllers.Permission.Insert)
	pvt.Put("/permissions/:id", middleware.CheckPermission("permissions", "update"), controllers.Permission.Update)
	pvt.Delete("/permissions/:id", middleware.CheckPermission("permissions", "delete"), controllers.Permission.Delete)

	pvt.Get("/assets", middleware.CheckPermission("assets", "read"), controllers.Asset.GetAll)
	pvt.Get("/assets/:id", middleware.CheckPermission("assets", "read"), controllers.Asset.Get)
	pvt.Get("/assets/:id/attributes", middleware.CheckPermission("assets", "read"), controllers.Asset.GetWithAttributes)
	pvt.Get("/assets/names", middleware.CheckPermission("assets", "read"), controllers.Asset.GetNames)
	pvt.Post("/assets", middleware.CheckPermission("assets", "create"), controllers.Asset.Insert)
	pvt.Put("/assets/:id", middleware.CheckPermission("assets", "update"), controllers.Asset.Update)
	pvt.Delete("/assets/:id", middleware.CheckPermission("assets", "delete"), controllers.Asset.Delete)

	pvt.Get("/attributes", middleware.CheckPermission("attributes", "read"), controllers.Attribute.GetAll)
	pvt.Get("/attributes/:id", middleware.CheckPermission("attributes", "read"), controllers.Attribute.Get)
	pvt.Post("/attributes", middleware.CheckPermission("attributes", "create"), controllers.Attribute.Insert)
	pvt.Put("/attributes/:id", middleware.CheckPermission("attributes", "update"), controllers.Attribute.Update)
	pvt.Delete("/attributes/:id", middleware.CheckPermission("attributes", "delete"), controllers.Attribute.Delete)
	// ...

	// Log start
	host := conf.Env.GetString("HOST")
	port := conf.Env.GetString("PORT")
	logger.Logger.Printf("Starting server on http://%s:%s\n", host, port)

	// Run server
	return app.Listen(conf.Env.GetString("HOST") + ":" + port)
}
