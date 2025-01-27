package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrderController struct {
	uc  interfaces.OrderUseCase
	log *logrus.Logger
}

func NewOrderController(usecase interfaces.OrderUseCase, logger *logrus.Logger) *OrderController {
	return &OrderController{
		uc:  usecase,
		log: logger,
	}
}

func (c *OrderController) CreateOrder(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateOrder)
	request.UserID = auth.ID
	if err := ctx.BodyParser(request); err != nil {
		c.log.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	if len(request.Items) == 0 {
		c.log.Warnf("Order items are empty")
		return fiber.NewError(fiber.StatusBadRequest, "Order items cannot be empty")
	}
	response, err := c.uc.CreateOrder(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to create order: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully created order", fiber.StatusCreated, nil, nil))
}

func (c *OrderController) GetOrdersByBuyer(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.SearchOrderRequest{
		UserID: auth.ID,
		Status: ctx.Query("status", ""),
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", 10),
	}
	response, total, err := c.uc.GetOrdersByBuyer(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to get orders: %+v", err)
		return err
	}
	paging := &model.PageMetadata{
		Page:      request.Page,
		Size:      request.Limit,
		TotalItem: total,
		TotalPage: int64(total) / int64(request.Limit),
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get orders", fiber.StatusOK, paging, nil))
}

func (c *OrderController) GetOrderByIdByBuyer(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	uuid := ctx.Params("order_uuid")
	request := &model.GetOrderDetails{
		UserID: auth.ID,
		OrderUUID: uuid,
	}
	response, err := c.uc.GetOrderByIdByBuyer(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to get order: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get order", fiber.StatusOK, nil, nil))
}

func (c *OrderController) CancelOrder(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.CancelOrderRequest{
		UserID: auth.ID,
		OrderUUID: ctx.Params("order_uuid"),
	}
	response, err := c.uc.CancelOrder(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to cancel order: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully cancel order", fiber.StatusOK, nil, nil))
}

func (c *OrderController) CheckoutOrder(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.CheckoutOrderRequest{
		UserID: auth.ID,
		OrderUUID: ctx.Params("order_uuid"),
	}
	response, err := c.uc.CheckoutOrder(ctx.UserContext(), request)
	if err != nil {
		c.log.Warnf("Failed to checkout order: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully checkout order", fiber.StatusOK, nil, nil))
}

func (c *OrderController) GetOrderByIdSeller(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := &model.GetOrderDetails{
		UserID: authID.ID,
		OrderUUID: ctx.Params("order_uuid"),
	}
	response, err := c.uc.GetOrderBySeller(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get order", fiber.StatusOK, nil, nil))
}

func (c *OrderController) GetOrdersBySeller(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := &model.SearchOrderRequestBySeller{
		UserID: authID.ID,
		Status: ctx.Query("status", ""),
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", 10),
	}
	response, total, err := c.uc.GetOrdersBySeller(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	paging := &model.PageMetadata{
		Page: request.Page,
		Size: request.Limit,
		TotalItem: total,
		TotalPage: int64(total) / int64(request.Limit),
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get orders", fiber.StatusOK, paging, nil))
}
