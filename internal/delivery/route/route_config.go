package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App *fiber.App
	AuthController *http.AuthController
}

func (c *RouteConfig) Setup() {
	SetupAuthRoute(c.App, c.AuthController)
}