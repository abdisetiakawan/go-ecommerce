package http

import (
	"time"

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
	authRes, err := c.uc.Register(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	ctx.Cookie(&fiber.Cookie{
        Name:     "jwt",
        Value:    authRes.AccessToken,
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Strict",
        MaxAge:   3600,
        Path:     "/",
    })
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(
        fiber.Map{"id": authRes.ID, "name": authRes.Name, "role": authRes.Role},
        "User registered successfully", fiber.StatusCreated, nil, nil,
    ))
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
	authRes, err := c.uc.Login(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	ctx.Cookie(&fiber.Cookie{
        Name:     "jwt",
        Value:    authRes.AccessToken,
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Strict",
        MaxAge:   3600,
        Path:     "/",
    })
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(
        fiber.Map{"id": authRes.ID, "name": authRes.Name, "role": authRes.Role},
        "User logged in successfully", fiber.StatusOK, nil, nil,
    ))
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

// Logout handles POST /users/logout endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request.
//
// Returns:
//
//   * 200 OK: model.WebResponse(true) if logout is successful.
//
// Errors:
//
//   * N/A
func (c *UserController) Logout(ctx *fiber.Ctx) error {
    ctx.Cookie(&fiber.Cookie{
        Name:     "jwt",
        Value:    "",
        HTTPOnly: true,
        Secure:   true,
        SameSite: "Strict",
        Expires:  time.Now().Add(-time.Hour),
        Path:     "/",
    })
    return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(true, "Logged out successfully", fiber.StatusOK, nil, nil))
}