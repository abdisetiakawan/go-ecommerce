package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/buyer"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/seller"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/user"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App *fiber.App
	AuthController *http.AuthController
	SellerController *seller.SellerController
	BuyerController *buyer.BuyerController
	UserController *user.UserController
	AuthMiddleware    fiber.Handler
}

func (c *RouteConfig) Setup() {
	SetupAuthRoute(c.App, c.AuthController)
	SetupUserRoute(c, c.App, c.UserController)
	SetupBuyerRoute(c, c.App, c.BuyerController)
	SetupSellerRoute(c, c.App, c.SellerController)
}