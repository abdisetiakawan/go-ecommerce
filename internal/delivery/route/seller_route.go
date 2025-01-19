package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/seller"
	"github.com/gofiber/fiber/v2"
)

func SetupSellerRoute(r *RouteConfig, app *fiber.App, sellerController *seller.SellerController) {
	app.Use(r.AuthMiddleware)
	app.Use(middleware.SellerOnly())
	app.Post("/api/seller/register", sellerController.RegisterStore)
	app.Post("/api/seller/products", sellerController.RegisterProduct)
}