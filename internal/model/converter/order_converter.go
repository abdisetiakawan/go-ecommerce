package converter

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	eventmodel "github.com/abdisetiakawan/go-ecommerce/internal/model/event_model"
)

func OrderToResponse(order *entity.Order) *model.OrderResponse {
    items := make([]model.OrderItemResponse, len(order.Items))

    for i, item := range order.Items {
        items[i] = model.OrderItemResponse{
            OrderItemUuid: item.OrderItemUUID,
            ProductName:   item.Product.ProductName,
            Price:         item.Product.Price,
            Quantity:  item.Quantity,
        }
    }

    return &model.OrderResponse{
        Username:   order.User.Username,
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
        CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
    }
}

func OrdersToResponse(orders *entity.Order) *model.ListOrderResponse {
    return &model.ListOrderResponse{
        OrderUUID: orders.OrderUUID,
        TotalPrice: orders.TotalPrice,
        Status: orders.Status,
        Date: orders.CreatedAt.Format("2006-01-02 15:04"),
    }
}

func OrderToResponseForSeller(order *entity.Order) *model.OrdersResponseForSeller {
    items := make([]model.OrderItemResponse, len(order.Items))

    for i, item := range order.Items {
        items[i] = model.OrderItemResponse{
            OrderItemUuid: item.OrderItemUUID,
            Quantity:  item.Quantity,
            Price:     item.Product.Price,
            ProductName: item.Product.ProductName,
        }
    }

    return &model.OrdersResponseForSeller{
        UserName:       order.User.Username,
        OrderUUID:  order.OrderUUID,
        TotalPrice: order.TotalPrice,
        Status:     order.Status,
        Items:      items,
        Payment: model.PaymentResponse{
            PaymentUUID:  order.Payment.PaymentUUID,
            PaymentMethod: order.Payment.Method,
            Status:        order.Payment.Status,
        },
        CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
    }
}

func CreateOrderToResponse(payment *eventmodel.PaymentMessage, shipping *eventmodel.ShippingMessage, order *entity.Order) *model.OrderResponse {
    items := make([]model.OrderItemResponse, len(order.Items))
    for i, item := range order.Items {
        items[i] = model.OrderItemResponse{
            OrderItemUuid: item.OrderItemUUID,
            Quantity:  item.Quantity,
        }
    }
    return &model.OrderResponse{
        Username:   order.User.Username,
        OrderUUID:  order.OrderUUID,
        TotalPrice: order.TotalPrice,
        Status:     order.Status,
        Items:      items,
        Shipping: model.ShippingResponse{
            ShippingUUID: shipping.ShippingUUID,
            Address:      shipping.Address,
            City:         shipping.City,
            Province:     shipping.Province,
            PostalCode:   shipping.PostalCode,
            Status:       shipping.Status,
        },
        Payment: model.PaymentResponse{
            PaymentUUID:  payment.PaymentUUID,
            PaymentMethod: payment.Method,
            Status:        payment.Status,
        },
        CreatedAt: order.CreatedAt.Format("2006-01-02 15:04:05"),
    }
}