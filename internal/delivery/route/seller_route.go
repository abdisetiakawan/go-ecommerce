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
	// Todo: Order Management:
	// Todo: Get Orders
	/*
	sellerGroup.Get("/orders", sellerController.GetOrders)
	*/
	// Todo: Get Order by ID
	/*
	sellerGroup.Get("/orders/:order_uuid", sellerController.GetOrderById)
	*/
	// Todo : Update shipping status:
	/*
	sellerGroup.Put("/orders/:order_uuid", sellerController.UpdateShippingStatus)
	*/
}