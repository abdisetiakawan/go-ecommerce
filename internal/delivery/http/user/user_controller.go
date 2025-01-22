package user

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	UseCase *usecase.UserUseCase
	Logger  *logrus.Logger
}

func NewUserController(usecase *usecase.UserUseCase, logger *logrus.Logger) *UserController {
	return &UserController{
		UseCase: usecase,
		Logger:  logger,
	}
}

func (c *UserController) CreateProfile(ctx *fiber.Ctx) error {
    authID := middleware.GetUser(ctx)
    request := new(model.CreateProfile)
    request.UserID = authID.ID

    if err := ctx.BodyParser(request); err != nil {
        c.Logger.Warnf("Failed to parse request body: %+v", err)
        return err
    }
    helper.TrimSpaces(request)
    response, err := c.UseCase.Create(ctx.UserContext(), request)
    if err != nil {
        c.Logger.Warnf("Failed to create profile: %+v", err)
        return err
    }

    return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully created profile", fiber.StatusCreated, nil, nil))
}

func (c *UserController) GetProfile(ctx *fiber.Ctx) error {
    authID := middleware.GetUser(ctx)
    response, err := c.UseCase.Get(ctx.UserContext(), authID.ID)
    if err != nil {
        c.Logger.Warnf("Failed to get profile: %+v", err)
        return err
    }

    return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get profile", fiber.StatusOK, nil, nil))
}

func (c *UserController) UpdateProfile(ctx *fiber.Ctx) error {
    authID := middleware.GetUser(ctx)
    request := new(model.UpdateProfile)
    request.UserID = authID.ID

    if err := ctx.BodyParser(request); err != nil {
        c.Logger.Warnf("Failed to parse request body: %+v", err)
        return err
    }
    helper.TrimSpaces(request)
    response, err := c.UseCase.Update(ctx.UserContext(), request)
    if err != nil {
        c.Logger.Warnf("Failed to update profile: %+v", err)
        return err
    }

    return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated profile", fiber.StatusOK, nil, nil))
}