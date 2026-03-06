package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	app "github.com/Dizao9/Fitness-Journal/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Apllication failed: %v", err)
	}

	log.Println("Application stopped safely")
}
