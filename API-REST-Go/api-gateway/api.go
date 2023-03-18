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
	// Create api router
	app := NewRouter()
	// Run server
	return app.Listen(conf.Env.GetString("HOST") + ":" + conf.Env.GetString("PORT"))
}

func NewRouter() *fiber.App {
	// Build controllers
	controllers.Build()

	// Setup Fiber api server
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	// Setup groups and middleware
	api := app.Group(conf.Conf.GetString("apiBasePath"), middleware.CORS())
	pub := api.Group("/public")
	pvt := api.Group("/private", middleware.CheckToken)
	pvtVer := pvt.Group("/verified", middleware.VerifyEmail)

	// Api status path
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

	// Exclusive path for serving all files (superadmin required)
	app.Get("/",
		middleware.CheckToken,
		middleware.CheckPermission("*", "*"),
		// Serve all files of project
		filesystem.New(filesystem.Config{
			Root:   http.Dir("./"),
			Browse: true,
		}),
	)

	pub.Post("/auth/login", controllers.User.Auth.Login)
	pub.Post("/auth/confirm-email/:token", controllers.User.Auth.ConfirmEmail)
	pub.Put("/auth/reset-password", controllers.User.Auth.ResetPassword)
	pvt.Post("/auth/resend-confirmation-email", controllers.User.Auth.ResendConfirmationEmail)

	pvt.Get("/users/me", controllers.User.Auth.Get)
	pvt.Get("/users/me/photo", controllers.User.Auth.GetPhoto)
	pvt.Get("/users/me/cv", controllers.User.Auth.GetCV)
	pvt.Put("/users/me", controllers.User.Auth.Update)
	pvt.Put("/users/me/change-password", controllers.User.Auth.ChangePassword)
	pvt.Put("/users/me/photo", controllers.User.Auth.UpdatePhoto)
	pvt.Put("/users/me/cv", controllers.User.Auth.UpdateCV)
	pvt.Delete("/users/me", controllers.User.Auth.Delete)

	pvtVer.Get("/users", middleware.CheckPermission("users", "read"), controllers.User.GetAll)
	pvtVer.Get("/users/:id", middleware.CheckPermission("users", "read"), controllers.User.Get)
	pvtVer.Get("/users/:id/photo", middleware.CheckPermission("users", "read"), controllers.User.GetPhoto)
	pvtVer.Get("/users/:id/cv", middleware.CheckPermission("users", "read"), controllers.User.GetCV)
	pub.Post("/users", controllers.User.Insert) // public registration
	pvtVer.Post("/users/:id/roles", middleware.CheckPermission("users", "assign"), controllers.User.UpdateRoles)
	pvtVer.Put("/users/:id", middleware.CheckPermission("users", "update"), controllers.User.Update)
	pvtVer.Put("/users/:id/photo", middleware.CheckPermission("users", "update"), controllers.User.UpdatePhoto)
	pvtVer.Put("/users/:id/cv", middleware.CheckPermission("users", "update"), controllers.User.UpdateCV)
	pvtVer.Put("/users/:id/ban", middleware.CheckPermission("users", "ban"), controllers.User.Ban)
	pvtVer.Put("/users/:id/unban", middleware.CheckPermission("users", "ban"), controllers.User.Unban)
	pvtVer.Put("/users/:id/restore", middleware.CheckPermission("users", "delete"), controllers.User.Restore)
	pvtVer.Delete("/users/:id", middleware.CheckPermission("users", "delete"), controllers.User.Delete)

	pvtVer.Get("/roles", middleware.CheckPermission("roles", "read"), controllers.Role.GetAll)
	pvtVer.Get("/roles/:id", middleware.CheckPermission("roles", "read"), controllers.Role.Get)
	pvtVer.Post("/roles", middleware.CheckPermission("roles", "create"), controllers.Role.Insert)
	pvtVer.Post("/roles/:id/permissions", middleware.CheckPermission("roles", "assign"), controllers.Role.UpdatePermissions)
	pvtVer.Put("/roles/:id", middleware.CheckPermission("roles", "update"), controllers.Role.Update)
	pvtVer.Delete("/roles/:id", middleware.CheckPermission("roles", "delete"), controllers.Role.Delete)

	pvtVer.Get("/permissions", middleware.CheckPermission("permissions", "read"), controllers.Permission.GetAll)
	pvtVer.Get("/permissions/:id", middleware.CheckPermission("permissions", "read"), controllers.Permission.Get)
	pvtVer.Post("/permissions", middleware.CheckPermission("permissions", "create"), controllers.Permission.Insert)
	pvtVer.Put("/permissions/:id", middleware.CheckPermission("permissions", "update"), controllers.Permission.Update)
	pvtVer.Delete("/permissions/:id", middleware.CheckPermission("permissions", "delete"), controllers.Permission.Delete)

	pvtVer.Get("/assets", middleware.CheckPermission("assets", "read"), controllers.Asset.GetAll)
	pvtVer.Get("/assets/:id", middleware.CheckPermission("assets", "read"), controllers.Asset.Get)
	pvtVer.Get("/assets/:id/attributes", middleware.CheckPermission("assets", "read"), controllers.Asset.GetWithAttributes)
	pvtVer.Get("/assets/names", middleware.CheckPermission("assets", "read"), controllers.Asset.GetNames)
	pvtVer.Post("/assets", middleware.CheckPermission("assets", "create"), controllers.Asset.Insert)
	pvtVer.Put("/assets/:id", middleware.CheckPermission("assets", "update"), controllers.Asset.Update)
	pvtVer.Delete("/assets/:id", middleware.CheckPermission("assets", "delete"), controllers.Asset.Delete)

	pvtVer.Get("/attributes", middleware.CheckPermission("attributes", "read"), controllers.Attribute.GetAll)
	pvtVer.Get("/attributes/:id", middleware.CheckPermission("attributes", "read"), controllers.Attribute.Get)
	pvtVer.Post("/attributes", middleware.CheckPermission("attributes", "create"), controllers.Attribute.Insert)
	pvtVer.Put("/attributes/:id", middleware.CheckPermission("attributes", "update"), controllers.Attribute.Update)
	pvtVer.Delete("/attributes/:id", middleware.CheckPermission("attributes", "delete"), controllers.Attribute.Delete)
	// ...

	// Log start
	host := conf.Env.GetString("HOST")
	port := conf.Env.GetString("PORT")
	logger.Logger.Printf("Starting server on http://%s:%s\n", host, port)

	return app
}
