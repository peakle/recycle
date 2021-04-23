package handlers

import (
	"time"

	_ "github.com/vrischmann/envconfig"
)

type Config struct {
	requestTimeout time.Time `envconfig:"default=2s"`
}
