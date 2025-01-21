package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/buyer"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupBuyerRoute(r *RouteConfig, app *fiber.App, buyerController *buyer.BuyerController) {
	buyerGroup := app.Group("/api/buyer", r.AuthMiddleware, middleware.BuyerOnly())
	buyerGroup.Get("/orders", buyerController.SearchOrders)
	buyerGroup.Post("/orders", buyerController.CreateOrder)
}