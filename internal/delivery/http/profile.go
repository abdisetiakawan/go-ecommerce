package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type ProfileController struct {
	uc  interfaces.ProfileUseCase
	log *logrus.Logger
}

func NewProfileController(usecase interfaces.ProfileUseCase, logger *logrus.Logger) *ProfileController {
	return &ProfileController{
		uc:  usecase,
		log: logger,
	}
}

func (c *ProfileController) CreateProfile(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.CreateProfile)
	request.UserID = authID.ID

	if err := ctx.BodyParser(request); err != nil {
		c.log.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.CreateProfile(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to create profile: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully created profile", fiber.StatusCreated, nil, nil))
}

func (c *ProfileController) GetProfile(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	response, err := c.uc.GetProfile(ctx.UserContext(), authID.ID)
	if err != nil {
		c.log.Warnf("Failed to get profile: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get profile", fiber.StatusOK, nil, nil))
}

func (c *ProfileController) UpdateProfile(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.UpdateProfile)
	request.UserID = authID.ID

	if err := ctx.BodyParser(request); err != nil {
		c.log.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.UpdateProfile(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to update profile: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated profile", fiber.StatusOK, nil, nil))
}
