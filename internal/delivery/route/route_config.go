package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/seller"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App *fiber.App
	AuthController *http.AuthController
	SellerController *seller.SellerController
	AuthMiddleware    fiber.Handler
}

func (c *RouteConfig) Setup() {
	SetupAuthRoute(c.App, c.AuthController)
	SetupSellerRoute(c, c.App, c.SellerController)
}