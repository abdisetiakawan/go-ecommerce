package interfaces

import (
	"time"

	evententity "github.com/abdisetiakawan/go-ecommerce/internal/entity/event_entity"
)

type OrderEventRepository interface {
	CreateOrderEvent(event *evententity.OrderEvent) error
	UpdateOrderEvent(event *evententity.OrderEvent) error
	GetPendingEvents() ([]evententity.OrderEvent, error)
	GetFailedEvents(duration time.Duration) ([]evententity.OrderEvent, error)
}
