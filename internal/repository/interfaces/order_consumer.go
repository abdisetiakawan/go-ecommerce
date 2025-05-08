package interfaces

import "context"

type OrderConsumer interface {
	ChangeOrderStatus(ctx context.Context) error
}
