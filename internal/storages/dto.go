package storages

import "time"

type Order struct {
	Id          string    `json:"id,omitempty"`
	Address     string    `json:"address,omitempty"`
	MaxSize     int       `json:"maxSize,omitempty"`
	CurrentSize int       `json:"currentSize,omitempty"`
	EventAt     time.Time `json:"eventAt,omitempty"`
}

type OrderUser struct {
	userId  string
	orderId string
}
