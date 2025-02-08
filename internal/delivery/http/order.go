package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
)

type OrderController struct {
	uc  interfaces.OrderUseCase
}

func NewOrderController(usecase interfaces.OrderUseCase) *OrderController {
	return &OrderController{
		uc:  usecase,
	}
}

// CreateOrder handles POST /orders endpoint. It creates a new order for user.
//
// Parameters:
//
//	* request body: model.CreateOrder
//
// Returns:
//
//	* 201 Created: model.OrderResponse
//
// Errors:
//
//	* 400 Bad Request: If order items are empty
//	* 500 Internal Server Error: If there's an error while creating order
func (c *OrderController) CreateOrder(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := new(model.CreateOrder)
	request.UserID = auth.ID
	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	if len(request.Items) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Order items cannot be empty")
	}
	response, err := c.uc.CreateOrder(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully created order", fiber.StatusCreated, nil, nil))
}

// GetOrdersByBuyer handles GET /orders endpoint for buyer.
//
// Parameters:
//
//	* ctx: fiber.Ctx - Context for the request, including query parameters for filtering orders by status, page, and limit.
//
// Returns:
//
//	* 200 OK: model.ListOrderResponse with pagination metadata if orders are retrieved successfully.
//
// Errors:
//
//	* Propagates error from use case layer if retrieval fails.
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

// GetOrderByIdByBuyer handles GET /orders/{order_uuid} endpoint for buyer.
//
// Parameters:
//
//	* ctx: fiber.Ctx - Context for the request, including the order UUID path parameter.
//
// Returns:
//
//	* 200 OK: model.OrderResponse if order is retrieved successfully.
//
// Errors:
//
//	* Propagates error from use case layer if retrieval fails.
func (c *OrderController) GetOrderByIdByBuyer(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	uuid := ctx.Params("order_uuid")
	request := &model.GetOrderDetails{
		UserID: auth.ID,
		OrderUUID: uuid,
	}
	response, err := c.uc.GetOrderByIdByBuyer(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get order", fiber.StatusOK, nil, nil))
}

// CancelOrder handles PATCH /orders/{order_uuid} endpoint for buyer.
//
// Parameters:
//
//	* ctx: fiber.Ctx - Context for the request, including the order UUID path parameter.
//
// Returns:
//
//	* 200 OK: model.OrderResponse if order is canceled successfully.
//
// Errors:
//
//	* Propagates error from use case layer if cancelation fails.
func (c *OrderController) CancelOrder(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.CancelOrderRequest{
		UserID: auth.ID,
		OrderUUID: ctx.Params("order_uuid"),
	}
	response, err := c.uc.CancelOrder(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully cancel order", fiber.StatusOK, nil, nil))
}

// CheckoutOrder handles PATCH /orders/{order_uuid}/checkout endpoint for buyer.
//
// Parameters:
//
//	* ctx: fiber.Ctx - Context for the request, including the order UUID path parameter.
//
// Returns:
//
//	* 200 OK: model.OrderResponse if order is checked out successfully.
//
// Errors:
//
//	* Propagates error from use case layer if checkout fails.
func (c *OrderController) CheckoutOrder(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	request := &model.CheckoutOrderRequest{
		UserID: auth.ID,
		OrderUUID: ctx.Params("order_uuid"),
	}
	response, err := c.uc.CheckoutOrder(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully checkout order", fiber.StatusOK, nil, nil))
}

// GetOrderByIdSeller handles GET /orders/{order_uuid} endpoint for seller.
//
// Parameters:
//
//	* ctx: fiber.Ctx - Context for the request, including the order UUID path parameter.
//
// Returns:
//
//	* 200 OK: model.OrderResponse if order is retrieved successfully.
//
// Errors:
//
//	* Propagates error from use case layer if retrieval fails.
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

// GetOrdersBySeller handles GET /orders endpoint for seller.
//
// Parameters:
//
//	* ctx: fiber.Ctx - Context for the request, including query parameters for filtering orders by status, page, and limit.
//
// Returns:
//
//	* 200 OK: model.OrdersResponseForSeller with pagination metadata if orders are retrieved successfully.
//
// Errors:
//
//	* Propagates error from use case layer if retrieval fails.
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
