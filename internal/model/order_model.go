package model

type CreateOrder struct {
	UserID          uint                   `json:"-"`
	Items           []OrderItemRequest     `json:"items" validate:"required,dive"`
	ShippingAddress ShippingAddressRequest `json:"shipping_address" validate:"required"`
	Payments        PaymentRequest         `json:"payments" validate:"required"`
}

type OrderItemRequest struct {
	ProductUUID string `json:"product_uuid" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required,gte=1"`
}

type ShippingAddressRequest struct {
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	Province   string `json:"province" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required,min=5,max=5"`
}

type PaymentRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=cash transfer"`
}

type OrderResponse struct {
	OrderUUID  string              `json:"order_uuid"`
	TotalPrice float64             `json:"total_price"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
	Shipping   ShippingResponse    `json:"shipping"`
	Payment    PaymentResponse     `json:"payment"`
}

type OrderItemResponse struct {
	OrderItemUuid string `json:"order_item_uuid"`
	Quantity  int  `json:"quantity"`
}

type ShippingResponse struct {
	ShippingUUID string `json:"shipping_uuid"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Province     string `json:"province"`
	PostalCode   string `json:"postal_code"`
	Status       string `json:"status"`
}

type PaymentResponse struct {
	PaymentUUID   string `json:"payment_uuid"`
	PaymentMethod string `json:"payment_method"`
	Status        string `json:"status"`
}