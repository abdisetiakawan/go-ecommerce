package middleware

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/gofiber/fiber/v2"
)

// Middleware to verify if the user is a seller
func SellerOnly() fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        auth := ctx.Locals("auth").(*model.Auth)
        if auth.Role != "seller" {
            return model.ErrForbidden
        }
        return ctx.Next()
    }
}

// Middleware to verify if the user is a buyer
func BuyerOnly() fiber.Handler {
    return func(ctx *fiber.Ctx) error {
        auth := ctx.Locals("auth").(*model.Auth)
        if auth.Role != "buyer" {
            return model.ErrForbidden
        }
        return ctx.Next()
    }
}