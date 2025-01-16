package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoute(app *fiber.App, authController *http.AuthController) {
	app.Post("/api/auth/register", authController.Register)
	app.Post("/api/auth/login", authController.Login)
}