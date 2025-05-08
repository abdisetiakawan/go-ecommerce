package interfaces

import "context"

type PaymentConsumer interface {
    CreatePayment(ctx context.Context) error
    CancelPayment(ctx context.Context) error 
    CheckoutPayment(ctx context.Context) error
}