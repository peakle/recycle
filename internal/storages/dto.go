package storages

type Order struct {
	Id          string `json:"order_id,omitempty"`
	Address     string `json:"address,omitempty"`
	MaxSize     int    `json:"max_size,omitempty"`
	CurrentSize int    `json:"current_size,omitempty"`
	EventAt     string `json:"event_at,omitempty"`
}

type OrderUser struct {
	userId  string
	orderId string
}
