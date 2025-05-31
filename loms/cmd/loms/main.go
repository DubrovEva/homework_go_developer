package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"route256/loms/internal/application"
)

func main() {
	fmt.Println("App `Loms` starting")

	app := application.NewApp()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v, initiating shutdown", sig)
		cancel()
	}()

	if err := app.Start(ctx); err != nil {
		if err := app.Close(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		log.Fatalf("Application error: %v", err)
	}

	<-ctx.Done()
	log.Println("Context cancelled, shutting down gracefully...")

	if err := app.Close(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Application shutdown complete")
}
