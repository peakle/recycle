package main

import (
	"context"

	"github.com/peakle/recycle/internal/server"
)

var (
	// Version - app release
	Version = "0"
	// CommitID - release's commit id
	CommitID = "0"
	commands = map[string]struct {
		Name        string
		Description string
		Action      func(context.Context) error
	}{
		"server": {
			Name:        "server",
			Description: "gamechart api server",
			Action:      server.StartServer,
		},
	}
)

func main() {

}
