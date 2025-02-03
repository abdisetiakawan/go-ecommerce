package model

type ShippingMessage struct {
	ShippingUUID string `json:"shipping_uuid"`
	OrderID uint   `json:"order_id"`
	Address string `json:"address"`
	City    string `json:"city"`
	Province string `json:"province"`
	PostalCode string `json:"postal_code"`
	Status string `json:"status"`
}