package middleware

import (
	"API-REST/services/conf"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORS() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOrigins: conf.Env.GetString("CORS_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	})
}
