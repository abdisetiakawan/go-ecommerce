package interfaces

import "context"

type ShippingConsumer interface {
    CreateShipping(ctx context.Context) error
    CancelShipping(ctx context.Context) error
}