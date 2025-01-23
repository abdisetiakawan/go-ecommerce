package buyer

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type BuyerController struct {
	UseCase *usecase.BuyerUseCase
	Logger  *logrus.Logger
}

func NewBuyerController(usecase *usecase.BuyerUseCase, logger *logrus.Logger) *BuyerController {
	return &BuyerController{
		UseCase: usecase,
		Logger:  logger,
	}
}

func (c *BuyerController) CreateOrder(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateOrder)
	request.UserID = auth.ID
	if err := ctx.BodyParser(request); err != nil {
		c.Logger.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	if len(request.Items) == 0 {
        c.Logger.Warnf("Order items are empty")
        return fiber.NewError(fiber.StatusBadRequest, "Order items cannot be empty")
    }
	response, err := c.UseCase.CreateOrder(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to create order: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully created order", fiber.StatusCreated, nil, nil))
}

func (c *BuyerController) SearchOrders(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.SearchOrderRequest{
		UserID: auth.ID,
		Status: ctx.Query("status", ""),
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", 10),
	}
	response, total, err := c.UseCase.GetOrders(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to get orders: %+v", err)
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

func (c *BuyerController) GetOrder(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	uuid := ctx.Params("order_uuid")
	request := &model.GetOrderDetails{
		UserID: auth.ID,
		OrderUUID: uuid,
	}
	response, err := c.UseCase.GetOrder(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to get order: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get order", fiber.StatusOK, nil, nil))
}