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
	sellerGroup.Get("/products/:product_uuid", sellerController.GetProductById)
	sellerGroup.Put("/products/:product_uuid", sellerController.UpdateProduct)
	sellerGroup.Delete("/products/:product_uuid", sellerController.DeleteProduct)
	sellerGroup.Get("/orders", sellerController.GetOrders)
	sellerGroup.Get("/orders/:order_uuid", sellerController.GetOrderById)
	sellerGroup.Patch("/orders/:order_uuid/shipping", sellerController.UpdateShippingStatus)
}