package auth

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	UseCase *usecase.AuthUseCase
	Logger *logrus.Logger
}

func NewAuthController(usecase *usecase.AuthUseCase, logger *logrus.Logger) *AuthController {
	return &AuthController{
		UseCase: usecase,
		Logger: logger,
	}
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUser)
	if err := ctx.BodyParser(request); err != nil {
		c.Logger.Warnf("Failed to parse request body : %+v", err)
		return err
	}
	helper.TrimSpaces(request, request.ConfirmPassword, request.Password)
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to register user : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered user", fiber.StatusCreated, nil, nil))
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUser)
	if err := ctx.BodyParser(request); err != nil {
		c.Logger.Warnf("Failed to parse request body : %+v", err)
		return err
	}
	response, err := c.UseCase.Login(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to login user : %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully logged in user", fiber.StatusOK, nil, nil))
}