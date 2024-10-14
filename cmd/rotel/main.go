package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	name := filepath.Base(os.Args[0])
	flags, err := NewArgs(name, os.Args[1:])
	if err == ErrHelp {
		os.Exit(0)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	// Create a context which cancels on CTRL+C
	ctx := HandleSignal()
	app, err := NewApp(ctx, flags.Name(), flags.Broker, flags.Credentials, flags.Qos, flags.Topic, flags.TTY)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	fmt.Println(flags)
	fmt.Println(app)

	// Run the application
	if err := app.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

// Handle signals - call cancel when interrupt received
func HandleSignal() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		cancel()
	}()
	return ctx
}
