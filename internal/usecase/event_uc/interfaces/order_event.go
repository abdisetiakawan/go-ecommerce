package interfaces

import (
	"context"

	evententity "github.com/abdisetiakawan/go-ecommerce/internal/entity/event_entity"
)

type OrderEventUseCase interface {
	ProcessOrderEvent(ctx context.Context, event *evententity.OrderEvent) error
	RetryFailedEvents(ctx context.Context) error
}