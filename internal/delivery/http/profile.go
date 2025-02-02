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

func (c *ProfileController) GetProfile(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	response, err := c.uc.GetProfile(ctx.UserContext(), authID.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get profile", fiber.StatusOK, nil, nil))
}

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
