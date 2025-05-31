package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"route256/cart/internal/application"
)

func main() {
	fmt.Println("App `cart` starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := application.NewApp()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("received signal: %s\n", sig)

		cancel()
	}()

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}
