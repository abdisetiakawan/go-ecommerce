package seller

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/delivery/http/middleware"
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
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to register store: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered store", fiber.StatusCreated, nil, nil))
}