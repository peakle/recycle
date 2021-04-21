package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
			Description: "recycle api server",
			Action:      server.StartServer,
		},
	}
)

func main() {
	commandName := flag.String("cmd", "", "command name")
	flag.Parse()

	fmt.Printf("%s-%s \n", Version, CommitID)

	var wg sync.WaitGroup
	var ctx, cancel = context.WithCancel(context.Background())
	var sigs = make(chan os.Signal)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	command, ok := commands[*commandName]
	if !ok {
		log.Printf("command not found: %s \n", *commandName)
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := command.Action(ctx)
		if err != nil {
			log.Printf("on main: %s", err)
		}
	}()

	wg.Wait()
}
