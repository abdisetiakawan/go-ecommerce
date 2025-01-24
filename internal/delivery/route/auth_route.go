package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/auth"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoute(app *fiber.App, authController *auth.AuthController) {
	app.Post("/api/auth/register", authController.Register)
	app.Post("/api/auth/login", authController.Login)
}