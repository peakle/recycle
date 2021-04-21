package handlers

import "time"

type Config struct {
	requestTimeout time.Time `envconfig:"default=2s"`
}
