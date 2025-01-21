package converter

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
)

func OrderToResponse(order *entity.Order) *model.OrderResponse {
    items := make([]model.OrderItemResponse, len(order.Items))

    for i, item := range order.Items {
        items[i] = model.OrderItemResponse{
            OrderItemUuid: item.OrderItemUUID,
            Quantity:  item.Quantity,
        }
    }

    return &model.OrderResponse{
        OrderUUID:  order.OrderUUID,
        TotalPrice: order.TotalPrice,
        Status:     order.Status,
        Items:      items,
        Shipping: model.ShippingResponse{
            ShippingUUID: order.Shipping.ShippingUUID,
            Address:      order.Shipping.Address,
            City:         order.Shipping.City,
            Province:     order.Shipping.Province,
            PostalCode:   order.Shipping.PostalCode,
            Status:       order.Shipping.Status,
        },
        Payment: model.PaymentResponse{
            PaymentUUID:  order.Payment.PaymentUUID,
            PaymentMethod: order.Payment.Method,
            Status:        order.Payment.Status,
        },
    }
}

func OrdersToResponse(orders *entity.Order) *model.ListOrderResponse {
    return &model.ListOrderResponse{
        OrderUUID: orders.OrderUUID,
        TotalPrice: orders.TotalPrice,
        Status: orders.Status,
    }
}