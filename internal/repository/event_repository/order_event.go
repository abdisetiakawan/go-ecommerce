package eventrepository

import (
	"time"

	evententity "github.com/abdisetiakawan/go-ecommerce/internal/entity/event_entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/event_repository/interfaces"
	"gorm.io/gorm"
)

type OrderEventRepository struct {
	DB *gorm.DB
}

func NewOrderEventRepository(db *gorm.DB) interfaces.OrderEventRepository {
	return &OrderEventRepository{DB: db}
}

func (r *OrderEventRepository) CreateOrderEvent(event *evententity.OrderEvent) error {
	return r.DB.Create(event).Error
}

func (r *OrderEventRepository) UpdateOrderEvent(event *evententity.OrderEvent) error {
	return r.DB.Save(event).Error
}

func (r *OrderEventRepository) GetPendingEvents() ([]evententity.OrderEvent, error) {
	var events []evententity.OrderEvent
	err := r.DB.Where("status = ?", "pending").Find(&events).Error
	return events, err
}

func (r *OrderEventRepository) GetFailedEvents(duration time.Duration) ([]evententity.OrderEvent, error) {
	var events []evententity.OrderEvent
	err := r.DB.Where("status = ? AND created_at > ?", "failed",
		time.Now().Add(-duration)).Find(&events).Error
	return events, err
}
