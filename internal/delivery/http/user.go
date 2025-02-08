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

// Register handles POST /users endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the request body for user registration.
//
// Returns:
//
//   * 201 Created: model.AuthResponse if user is registered successfully.
//
// Errors:
//
//   * Propagates error from use case layer if registration fails.
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

// Login handles POST /users/login endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the request body for user login.
//
// Returns:
//
//   * 200 OK: model.AuthResponse if user is logged in successfully.
//
// Errors:
//
//   * Propagates error from use case layer if login fails.
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

// ChangePassword handles PATCH /users/change_password endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the authenticated user and request body for changing password.
//
// Returns:
//
//   * 200 OK: model.WebResponse(true) if password is changed successfully.
//
// Errors:
//
//   * Propagates error from use case layer if change password fails.
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