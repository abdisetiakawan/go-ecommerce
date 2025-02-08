package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
)

type ProfileController struct {
	uc  interfaces.ProfileUseCase
}

func NewProfileController(usecase interfaces.ProfileUseCase) *ProfileController {
	return &ProfileController{
		uc:  usecase,
	}
}

// CreateProfile handles POST /profiles endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the authenticated user and request body for profile creation.
//
// Returns:
//
//   * 201 Created: model.ProfileResponse if profile is created successfully.
//
// Errors:
//
//   * Propagates error from use case layer if creation fails.
func (c *ProfileController) CreateProfile(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.CreateProfile)
	request.UserID = authID.ID

	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.CreateProfile(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully created profile", fiber.StatusCreated, nil, nil))
}

// GetProfile handles GET /profiles endpoint for getting authenticated user's profile.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the authenticated user.
//
// Returns:
//
//   * 200 OK: model.ProfileResponse if profile is retrieved successfully.
//
// Errors:
//
//   * Propagates error from use case layer if retrieval fails.
func (c *ProfileController) GetProfile(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	response, err := c.uc.GetProfile(ctx.UserContext(), authID.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get profile", fiber.StatusOK, nil, nil))
}

/*************  ✨ Codeium Command ⭐  *************/
// UpdateProfile handles PATCH /profiles endpoint for updating authenticated user's profile.
//
// Parameters:
//

//   * ctx: fiber.Ctx - Context for the request, including the authenticated user and request body for profile update.
//
// Returns:
//
//   * 200 OK: model.ProfileResponse if profile is updated successfully.
//
// Errors:
//
//   * Propagates error from use case layer if update fails.

func (c *ProfileController) UpdateProfile(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.UpdateProfile)
	request.UserID = authID.ID

	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.UpdateProfile(ctx.UserContext(), request)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated profile", fiber.StatusOK, nil, nil))
}
