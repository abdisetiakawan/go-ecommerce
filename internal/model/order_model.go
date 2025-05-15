package model

type CreateOrder struct {
	UserID          uint                   `json:"-"`
	Items           []OrderItemRequest     `json:"items" validate:"required,dive"`
	ShippingAddress ShippingAddressRequest `json:"shipping_address" validate:"required"`
	Payments        PaymentRequest         `json:"payments" validate:"required"`
}

type OrderItemRequest struct {
	ProductUUID string `json:"product_uuid" validate:"required,uuid"`
	Quantity    int    `json:"quantity" validate:"required,gte=1"`
}

type ShippingAddressRequest struct {
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	Province   string `json:"province" validate:"required"`
	PostalCode string `json:"postal_code" validate:"required,len=5,numeric"`
}

type PaymentRequest struct {
	PaymentMethod string `json:"payment_method" validate:"required,oneof=cash transfer"`
}

type OrderResponse struct {
	Username   string              `json:"user_name,omitempty"`
	OrderUUID  string              `json:"order_uuid"`
	TotalPrice float64             `json:"total_price"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
	Shipping   ShippingResponse    `json:"shipping"`
	Payment    PaymentResponse     `json:"payment"`
	CreatedAt  string              `json:"created_at"`
}

type OrderItemResponse struct {
	OrderItemUuid string  `json:"order_item_uuid"`
	ProductName   string  `json:"product_name"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
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

type SearchOrderRequest struct {
	UserID uint   `json:"-"`
	Status string `json:"-" validate:"omitempty,oneof=pending processed delivered cancelled"`
	Page   int    `json:"-"`
	Limit  int    `json:"-"`
}

type ListOrderResponse struct {
	OrderUUID  string  `json:"order_uuid"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
	Date       string  `json:"date"`
}

type GetOrderDetails struct {
	UserID    uint   `json:"-"`
	OrderUUID string `json:"-" validate:"required,uuid"`
}

type SearchOrderRequestBySeller struct {
	Status   string `json:"-" validate:"omitempty,oneof=pending processed delivered cancelled"`
	SortDate string `json:"-" validate:"omitempty,oneof=asc desc"`
	UserID   uint   `json:"-"`
	StoreID  uint   `json:"-"`
	Page     int    `json:"-"`
	Limit    int    `json:"-"`
}

type OrdersResponseForSeller struct {
	UserName   string              `json:"user_name"`
	OrderUUID  string              `json:"order_uuid"`
	TotalPrice float64             `json:"total_price"`
	Status     string              `json:"status"`
	Items      []OrderItemResponse `json:"items"`
	Payment    PaymentResponse     `json:"payment"`
	CreatedAt  string              `json:"created_at"`
}

type CancelOrderRequest struct {
	OrderUUID string `json:"-" validate:"required,uuid"`
	UserID    uint   `json:"-"`
}

type CheckoutOrderRequest struct {
	OrderUUID string `json:"-" validate:"required,uuid"`
	UserID    uint   `json:"-"`
}

type UpdateShippingStatusRequest struct {
	OrderUUID string `json:"-" validate:"required,uuid"`
	UserID    uint   `json:"-"`
	Status    string `json:"status" validate:"required,oneof=shipped delivered"`
}