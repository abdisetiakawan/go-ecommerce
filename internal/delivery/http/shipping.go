package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
)

type ShippingController struct {
	uc  interfaces.ShippingUseCase
}

func NewShippingController(usecase interfaces.ShippingUseCase) *ShippingController {
	return &ShippingController{
		uc:  usecase,
	}
}

func (c *ShippingController) UpdateShippingStatus(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := &model.UpdateShippingStatusRequest{
		OrderUUID: ctx.Params("order_uuid"),
		UserID:    authID.ID,
	}
	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	response, err := c.uc.UpdateShippingStatus(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated shipping status", fiber.StatusOK, nil, nil))
}