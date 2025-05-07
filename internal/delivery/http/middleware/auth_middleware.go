package middleware

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

func NewAuth(v *viper.Viper) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenStr := c.Cookies("jwt")
        if tokenStr == "" {
            return fiber.ErrUnauthorized
        }
        secret := v.GetString("credentials.accesssecret")
        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
        if err != nil || !token.Valid {
            return fiber.ErrUnauthorized
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return fiber.ErrUnauthorized
        }
        auth := &model.Auth{
            ID:    uint(claims["id"].(float64)),
            Email: claims["email"].(string),
            Role:  claims["role"].(string),
        }
        c.Locals("auth", auth)
        return c.Next()
    }
}

func GetUser(c *fiber.Ctx) *model.Auth {
    return c.Locals("auth").(*model.Auth)
}