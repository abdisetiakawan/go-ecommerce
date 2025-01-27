package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type StoreController struct {
	uc  interfaces.StoreUseCase
	log *logrus.Logger
}

func NewStoreController(usecase interfaces.StoreUseCase, logger *logrus.Logger) *StoreController {
	return &StoreController{
		uc:  usecase,
		log: logger,
	}
}

func (c *StoreController) RegisterStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.RegisterStore)
	request.ID = authID.ID
	if err := ctx.BodyParser(request); err != nil {
		c.log.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.RegisterStore(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to register store: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered store", fiber.StatusCreated, nil, nil))
}

func (c *StoreController) GetStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	response, err := c.uc.GetStore(ctx.UserContext(), authID.ID)
	if err != nil {
		c.log.Warnf("Failed to get store: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get store", fiber.StatusOK, nil, nil))
}

func (c *StoreController) UpdateStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.UpdateStore)
	request.ID = authID.ID
	if err := ctx.BodyParser(request); err != nil {
		c.log.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.UpdateStore(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to update store: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated store", fiber.StatusOK, nil, nil))
}
