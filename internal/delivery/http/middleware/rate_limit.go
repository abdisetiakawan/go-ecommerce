package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type RateLimiterType string

const (
    Auth   RateLimiterType = "auth"
    Seller RateLimiterType = "seller"
    Buyer  RateLimiterType = "buyer"
    User   RateLimiterType = "user"
)

type RateLimiterConfig struct {
    Type       RateLimiterType
    Max        int
    Expiration time.Duration
}

func createLimiterConfig(cfg *RateLimiterConfig) limiter.Config {
    return limiter.Config{
        Max:        cfg.Max,
        Expiration: cfg.Expiration,
        KeyGenerator: func(c *fiber.Ctx) string {
            if cfg.Type == Auth {
                return c.IP()
            }
            auth := GetUser(c)
            return c.IP() + auth.Email
        },
        LimitReached: func(c *fiber.Ctx) error {
            message := "rate limit exceeded"
            if cfg.Type != Auth {
                message = string(cfg.Type) + " " + message
            }
            return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
                "error": message + ", please try again later.",
            })
        },
    }
}

func NewDynamicRateLimiter(maxRequests int, expiration time.Duration) fiber.Handler {
    config := RateLimiterConfig{
        Type:       Auth,
        Max:        maxRequests,
        Expiration: expiration,
    }
    return limiter.New(createLimiterConfig(&config))
}

func NewSellerRateLimiter(maxRequests int, expiration time.Duration) fiber.Handler {
    config := RateLimiterConfig{
        Type:       Seller,
        Max:        maxRequests,
        Expiration: expiration,
    }
    return limiter.New(createLimiterConfig(&config))
}

func NewBuyerRateLimiter(maxRequests int, expiration time.Duration) fiber.Handler {
    config := RateLimiterConfig{
        Type:       Buyer,
        Max:        maxRequests,
        Expiration: expiration,
    }
    return limiter.New(createLimiterConfig(&config))
}

func NewUserRateLimiter(maxRequests int, expiration time.Duration) fiber.Handler {
    config := RateLimiterConfig{
        Type:       User,
        Max:        maxRequests,
        Expiration: expiration,
    }
    return limiter.New(createLimiterConfig(&config))
}