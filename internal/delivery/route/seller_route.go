package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/seller"
	"github.com/gofiber/fiber/v2"
)

func SetupSellerRoute(app *fiber.App, sellerController *seller.StoreController) {
	app.Post("/api/seller/register", sellerController.RegisterStore)
}