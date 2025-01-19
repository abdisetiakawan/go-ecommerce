package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/user"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoute(r *RouteConfig, app *fiber.App, userController *user.UserController) {
	app.Use(r.AuthMiddleware)
	app.Post("/api/user/profile", userController.CreateProfile)
}