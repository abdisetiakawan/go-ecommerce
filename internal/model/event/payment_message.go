package event

type PaymentMessage struct {
	PaymentUUID string `json:"payment_uuid"`
	OrderID uint   `json:"order_id"`
	Status string  `json:"status"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
}