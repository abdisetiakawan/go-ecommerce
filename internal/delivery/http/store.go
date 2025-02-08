package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
)

type StoreController struct {
	uc  interfaces.StoreUseCase
}

func NewStoreController(usecase interfaces.StoreUseCase) *StoreController {
	return &StoreController{
		uc:  usecase,
	}
}

// RegisterStore handles POST /stores endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the authenticated user and request body for store registration.
//
// Returns:
//
//   * 201 Created: model.StoreResponse if store is registered successfully.
//
// Errors:
//
//   * Propagates error from use case layer if registration fails.
func (c *StoreController) RegisterStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.RegisterStore)
	request.ID = authID.ID
	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.RegisterStore(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered store", fiber.StatusCreated, nil, nil))
}

// GetStore handles GET /stores endpoint for getting authenticated user's store.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the authenticated user.
//
// Returns:
//
//   * 200 OK: model.StoreResponse if store is retrieved successfully.
//
// Errors:
//
//   * Propagates error from use case layer if retrieval fails.
func (c *StoreController) GetStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	response, err := c.uc.GetStore(ctx.UserContext(), authID.ID)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get store", fiber.StatusOK, nil, nil))
}

// UpdateStore handles PATCH /stores endpoint for updating authenticated user's store.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the authenticated user and request body for updating store.
//
// Returns:
//
//   * 200 OK: model.StoreResponse if store is updated successfully.
//
// Errors:
//
//   * Propagates error from use case layer if update fails.
func (c *StoreController) UpdateStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.UpdateStore)
	request.ID = authID.ID
	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.UpdateStore(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated store", fiber.StatusOK, nil, nil))
}
