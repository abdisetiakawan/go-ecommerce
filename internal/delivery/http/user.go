package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	uc interfaces.UserUseCase
}

func NewUserController(usecase interfaces.UserUseCase) *UserController {
	return &UserController{
		uc: usecase,
	}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUser)
	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	helper.TrimSpaces(request, request.ConfirmPassword, request.Password)
	response, err := c.uc.Register(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered user", fiber.StatusCreated, nil, nil))
}

func (c *UserController) Login(ctx *fiber.Ctx) error {
	request := new(model.LoginUser)
	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	response, err := c.uc.Login(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully logged in user", fiber.StatusOK, nil, nil))
}

func (c *UserController) ChangePassword(ctx *fiber.Ctx) error {
    authID := middleware.GetUser(ctx)
    request := new(model.ChangePassword)
    request.UserID = authID.ID

    if err := ctx.BodyParser(request); err != nil {
        return err
    }
    helper.TrimSpaces(request, request.Password, request.ConfirmPassword, request.OldPassword)
    if err := c.uc.ChangePassword(ctx.UserContext(), request); err != nil {
        return err
    }

    return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(true, "Successfully change password", fiber.StatusOK, nil, nil))
}