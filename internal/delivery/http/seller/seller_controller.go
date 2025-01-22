package seller

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SellerController struct {
	UseCase *usecase.StoreUseCase
	Logger  *logrus.Logger
}

func NewSellerController(usecase *usecase.StoreUseCase, logger *logrus.Logger) *SellerController {
	return &SellerController{
		UseCase: usecase,
		Logger:  logger,
	}
}

func (c *SellerController) RegisterStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.RegisterStore)
	request.ID = authID.ID
	if err := ctx.BodyParser(request); err != nil {
		c.Logger.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to register store: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered store", fiber.StatusCreated, nil, nil))
}

func (c *SellerController) RegisterProduct(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.RegisterProduct)
	request.AuthID = authID.ID
	if err := ctx.BodyParser(request); err != nil {
		c.Logger.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.UseCase.CreateProduct(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to register product: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered product", fiber.StatusCreated, nil, nil))
}

func (c *SellerController) GetStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	response, err := c.UseCase.GetStore(ctx.UserContext(), authID.ID)
	if err != nil {
		c.Logger.Warnf("Failed to get store: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get store", fiber.StatusOK, nil, nil))
}

func (c *SellerController) UpdateStore(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.UpdateStore)
	request.ID = authID.ID
	if err := ctx.BodyParser(request); err != nil {
		c.Logger.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to update store: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated store", fiber.StatusOK, nil, nil))
}

func (c *SellerController) GetProducts(ctx *fiber.Ctx) error {
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
	response, total, err := c.UseCase.GetProducts(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to get products: %+v", err)
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

func (c *SellerController) GetProductById(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := &model.GetProductRequest{
		UserID: authID.ID,
		ProductUUID: ctx.Params("product_uuid"),
	}
	response, err := c.UseCase.GetProduct(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to get product: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully get product", fiber.StatusOK, nil, nil))
}

func (c *SellerController) UpdateProduct(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := new(model.UpdateProduct)
	request.UserID = authID.ID
	request.ProductUUID = ctx.Params("product_uuid")
	if err := ctx.BodyParser(request); err != nil {
		c.Logger.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	helper.TrimSpaces(request)
	response, err := c.UseCase.UpdateProduct(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to update product: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(model.NewWebResponse(response, "Successfully updated product", fiber.StatusOK, nil, nil))
}

func (c *SellerController) DeleteProduct(ctx *fiber.Ctx) error {
	authID := middleware.GetUser(ctx)
	request := &model.DeleteProductRequest{
		UserID: authID.ID,
		ProductUUID: ctx.Params("product_uuid"),
	}
	if err := c.UseCase.DeleteProduct(ctx.UserContext(), request); err != nil {
		c.Logger.Warnf("Failed to delete product: %+v", err)
		return err
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}