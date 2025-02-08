package eventmodel

type OrderMessage struct {
	OrderID    uint    `json:"order_id"`
	OrderUUID  string  `json:"order_uuid"`
	UserID     uint    `json:"user_id"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"total_price"`
}