package route

import (
	"time"

	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http"
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App               *fiber.App
	UserController    *http.UserController
	ProfileController *http.ProfileController
	OrderController   *http.OrderController
	StoreController   *http.StoreController
	ProductController *http.ProductController
	ShippingController *http.ShippingController
	AuthMiddleware    fiber.Handler
}

func (rc *RouteConfig) Setup() {
	rc.setupAuthRoutes()
	rc.setupUserRoutes()
	rc.setupBuyerRoutes()
	rc.setupSellerRoutes()
}

func (rc *RouteConfig) setupAuthRoutes() {
	authRateLimiter := middleware.NewDynamicRateLimiter(5, 2*time.Minute) 
	authGroup := rc.App.Group("/api/auth", authRateLimiter)
	{
		authGroup.Post("/register", rc.UserController.Register)
		authGroup.Post("/login", rc.UserController.Login)
	}
}

func (rc *RouteConfig) setupUserRoutes() {
	userRateLimiter := middleware.NewUserRateLimiter(40, time.Minute)
	userGroup := rc.App.Group("/api/user", rc.AuthMiddleware, userRateLimiter)
	{
		userGroup.Post("/profile", rc.ProfileController.CreateProfile)
		userGroup.Get("/profile", rc.ProfileController.GetProfile)
		userGroup.Put("/profile", rc.ProfileController.UpdateProfile)
		userGroup.Patch("/password", rc.UserController.ChangePassword)
	}
}

func (rc *RouteConfig) setupBuyerRoutes() {
	buyerRateLimiter := middleware.NewBuyerRateLimiter(30, time.Minute)
	buyerGroup := rc.App.Group("/api/buyer", rc.AuthMiddleware, middleware.BuyerOnly(), buyerRateLimiter)
	{
		// Order Routes
		orderGroup := buyerGroup.Group("/orders")
		{
			orderGroup.Get("", rc.OrderController.GetOrdersByBuyer)
			orderGroup.Get("/:order_uuid", rc.OrderController.GetOrderByIdByBuyer)
			orderGroup.Post("", rc.OrderController.CreateOrder)
			orderGroup.Patch("/:order_uuid/cancel", rc.OrderController.CancelOrder)
			orderGroup.Patch("/:order_uuid/checkout", rc.OrderController.CheckoutOrder)
		}
	}
}

func (rc *RouteConfig) setupSellerRoutes() {
	sellerRateLimiter := middleware.NewSellerRateLimiter(50, time.Minute)
	sellerGroup := rc.App.Group("/api/seller", rc.AuthMiddleware, middleware.SellerOnly(), sellerRateLimiter)
	{
		// Store Routes
		storeGroup := sellerGroup.Group("/store")
		{
			storeGroup.Post("", rc.StoreController.RegisterStore)
			storeGroup.Get("", rc.StoreController.GetStore)
			storeGroup.Put("", rc.StoreController.UpdateStore)
		}

		// Product Routes
		productGroup := sellerGroup.Group("/products")
		{
			productGroup.Post("", rc.ProductController.RegisterProduct)
			productGroup.Get("", rc.ProductController.GetProducts)
			productGroup.Get("/:product_uuid", rc.ProductController.GetProductById)
			productGroup.Put("/:product_uuid", rc.ProductController.UpdateProduct)
			productGroup.Delete("/:product_uuid", rc.ProductController.DeleteProduct)
		}

		// Order Routes
		orderGroup := sellerGroup.Group("/orders")
		{
			orderGroup.Get("", rc.OrderController.GetOrdersBySeller)
			orderGroup.Get("/:order_uuid", rc.OrderController.GetOrderByIdSeller)
			orderGroup.Patch("/:order_uuid/shipping", rc.ShippingController.UpdateShippingStatus)
		}
	}
}
