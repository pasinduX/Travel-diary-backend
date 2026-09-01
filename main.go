package main

import (
	"log"

	"travel-diary-backend/internal/app"
	"travel-diary-backend/internal/config"
)

func main() {
	cfg := config.Load()

	fiberApp, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := fiberApp.Listen(cfg.ServerAddress()); err != nil {
		log.Fatal(err)
	}
}
