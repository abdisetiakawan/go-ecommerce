package http

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	uc  interfaces.ProductUseCase
}

func NewProductController(usecase interfaces.ProductUseCase) *ProductController {
	return &ProductController{
		uc: usecase,
	}
}

// GetProducts handles GET /products endpoint for buyer and seller.
//
// Parameters:
//
//	* ctx: fiber.Ctx - Context for the request, including query parameters for filtering products by search, category, price range, page, and limit.
//
// Returns:
//
//	* 200 OK: model.ListProductResponse with pagination metadata if products are retrieved successfully.
//
// Errors:
//
//	* Propagates error from use case layer if retrieval fails.
func (c *ProductController) GetProducts(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := &model.GetProductsRequest{
		UserID: authID.ID,
		Search: ctx.Query("search", ""),
		Category: ctx.Query("category", ""),
		PriceMin: ctx.QueryFloat("price_min"),
		PriceMax: ctx.QueryFloat("price_max"),
		Page:   ctx.QueryInt("page", 1),
		Limit:  ctx.QueryInt("limit", 10),
	}
	response, total, err := c.uc.GetProducts(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	paging := &model.PageMetadata{
		Page: request.Page,
		Size: request.Limit,
		TotalItem: total,
		TotalPage: int64(total) / int64(request.Limit),
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get products", fiber.StatusOK, paging, nil))
}


// GetProductById handles GET /products/{product_uuid} endpoint for buyer and seller.
//
// Parameters:
//
//	* ctx: fiber.Ctx - Context for the request, including the product UUID path parameter.
//
// Returns:
//
//	* 200 OK: model.ProductResponse if product is retrieved successfully.
//
// Errors:
//
//	* Propagates error from use case layer if retrieval fails.
func (c *ProductController) GetProductById(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := &model.GetProductRequest{
		UserID: authID.ID,
		ProductUUID: ctx.Params("product_uuid"),
	}
	response, err := c.uc.GetProductById(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get product", fiber.StatusOK, nil, nil))
}

// UpdateProduct handles PATCH /products/{product_uuid} endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the product UUID path parameter.
//
// Returns:
//
//   * 200 OK: model.ProductResponse if product is updated successfully.
//
// Errors:
//
//   * Propagates error from use case layer if update fails.

func (c *ProductController) UpdateProduct(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.UpdateProduct)
	request.UserID = authID.ID
	request.ProductUUID = ctx.Params("product_uuid")
	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.UpdateProduct(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated product", fiber.StatusOK, nil, nil))
}



// DeleteProduct handles DELETE /products/{product_uuid} endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the product UUID path parameter.
//
// Returns:
//
//   * 204 No Content: if the product is deleted successfully.
//
// Errors:
//
//   * Propagates error from use case layer if deletion fails.

func (c *ProductController) DeleteProduct(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := &model.DeleteProductRequest{
		UserID: authID.ID,
		ProductUUID: ctx.Params("product_uuid"),
	}
	if err := c.uc.DeleteProduct(ctx.UserContext(), request); err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}
// RegisterProduct handles POST /products endpoint.
//
// Parameters:
//
//   * ctx: fiber.Ctx - Context for the request, including the authenticated user and request body for product registration.
//
// Returns:
//
//   * 201 Created: model.ProductResponse if product is registered successfully.
//
// Errors:
//
//   * Propagates error from use case layer if registration fails.
func (c *ProductController) RegisterProduct(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.RegisterProduct)
	request.AuthID = authID.ID
	if err := ctx.BodyParser(request); err != nil {
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.uc.CreateProduct(ctx.UserContext(), request)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered product", fiber.StatusCreated, nil, nil))
}
