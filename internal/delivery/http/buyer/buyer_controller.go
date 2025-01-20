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
	response, err := c.UseCase.CreateOrder(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to create order: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully created order", fiber.StatusCreated, nil, nil))
}