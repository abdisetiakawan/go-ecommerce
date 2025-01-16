package config

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
    var app = fiber.New(fiber.Config{
        AppName:      config.GetString("app.name"),
        ErrorHandler: NewErrorHandler(),
        Prefork:      config.GetBool("web.prefork"),
    })

    return app
}

func NewErrorHandler() fiber.ErrorHandler {
    return func(ctx *fiber.Ctx, err error) error {
        code := fiber.StatusInternalServerError
        message := err.Error()

        switch e := err.(type) {
        case *model.ApiError:
            code = e.StatusCode
            message = e.Message
        case *fiber.Error:
            code = e.Code
            message = e.Message
        }
        return ctx.Status(code).JSON(fiber.Map{
            "errors": message,
        })
    }
}