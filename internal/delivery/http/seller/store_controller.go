package seller

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type StoreController struct {
	UseCase *usecase.StoreUseCase
	Logger  *logrus.Logger
}

func NewStoreController(usecase *usecase.StoreUseCase, logger *logrus.Logger) *StoreController {
	return &StoreController{
		UseCase: usecase,
		Logger:  logger,
	}
}

func (c *StoreController) RegisterStore(ctx *fiber.Ctx) error {
	request := new(model.RegisterStore)
	if err := ctx.BodyParser(request); err != nil {
		c.Logger.Warnf("Failed to parse request body: %+v", err)
		return err
	}
	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Logger.Warnf("Failed to register store: %+v", err)
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(model.NewWebResponse(response, "Successfully registered store", fiber.StatusCreated, nil))
}