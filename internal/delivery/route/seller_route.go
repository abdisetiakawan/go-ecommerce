package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/seller"
	"github.com/gofiber/fiber/v2"
)

func SetupSellerRoute(r *RouteConfig, app *fiber.App, sellerController *seller.SellerController) {
	sellerGroup := app.Group("/api/seller", r.AuthMiddleware, middleware.SellerOnly())
	sellerGroup.Post("/store", sellerController.RegisterStore)
	sellerGroup.Get("/store", sellerController.GetStore)
	sellerGroup.Put("/store", sellerController.UpdateStore)
	sellerGroup.Post("/products", sellerController.RegisterProduct)
	sellerGroup.Get("/products", sellerController.GetProducts)
}