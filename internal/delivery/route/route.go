package route

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/auth"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/buyer"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/seller"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/user"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	AuthController    *auth.AuthController
	SellerController  *seller.SellerController
	BuyerController   *buyer.BuyerController
	UserController    *user.UserController
	AuthMiddleware    fiber.Handler
}

func (rc *RouteConfig) Setup() {
	rc.setupAuthRoutes()
	rc.setupUserRoutes()
	rc.setupBuyerRoutes()
	rc.setupSellerRoutes()
}

func (rc *RouteConfig) setupAuthRoutes() {
	authGroup := rc.App.Group("/api/auth")
	{
		authGroup.Post("/register", rc.AuthController.Register)
		authGroup.Post("/login", rc.AuthController.Login)
	}
}

func (rc *RouteConfig) setupUserRoutes() {
	userGroup := rc.App.Group("/api/user", rc.AuthMiddleware)
	{
		userGroup.Post("/profile", rc.UserController.CreateProfile)
		userGroup.Get("/profile", rc.UserController.GetProfile)
		userGroup.Put("/profile", rc.UserController.UpdateProfile)
		userGroup.Patch("/password", rc.UserController.ChangePassword)
	}
}

func (rc *RouteConfig) setupBuyerRoutes() {
	buyerGroup := rc.App.Group("/api/buyer", rc.AuthMiddleware, middleware.BuyerOnly())
	{
		// Order Routes
		orderGroup := buyerGroup.Group("/orders")
		{
			orderGroup.Get("", rc.BuyerController.SearchOrders)
			orderGroup.Get("/:order_uuid", rc.BuyerController.GetOrder)
			orderGroup.Post("", rc.BuyerController.CreateOrder)
			orderGroup.Patch("/:order_uuid/cancel", rc.BuyerController.CancelOrder)
			orderGroup.Patch("/:order_uuid/checkout", rc.BuyerController.CheckoutOrder)
		}
	}
}

func (rc *RouteConfig) setupSellerRoutes() {
	sellerGroup := rc.App.Group("/api/seller", rc.AuthMiddleware, middleware.SellerOnly())
	{
		// Store Routes
		storeGroup := sellerGroup.Group("/store")
		{
			storeGroup.Post("", rc.SellerController.RegisterStore)
			storeGroup.Get("", rc.SellerController.GetStore)
			storeGroup.Put("", rc.SellerController.UpdateStore)
		}

		// Product Routes
		productGroup := sellerGroup.Group("/products")
		{
			productGroup.Post("", rc.SellerController.RegisterProduct)
			productGroup.Get("", rc.SellerController.GetProducts)
			productGroup.Get("/:product_uuid", rc.SellerController.GetProductById)
			productGroup.Put("/:product_uuid", rc.SellerController.UpdateProduct)
			productGroup.Delete("/:product_uuid", rc.SellerController.DeleteProduct)
		}

		// Order Routes
		orderGroup := sellerGroup.Group("/orders")
		{
			orderGroup.Get("", rc.SellerController.GetOrders)
			orderGroup.Get("/:order_uuid", rc.SellerController.GetOrderById)
			orderGroup.Patch("/:order_uuid/shipping", rc.SellerController.UpdateShippingStatus)
		}
	}
}
