package main

import (
	"log"
	"net/http"
	"sync"

	"travel-diary-backend/internal/app"
	"travel-diary-backend/internal/config"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
)

var (
	serverlessApp *fiber.App
	serverlessErr  error
	once           sync.Once
)

func getApp() (*fiber.App, error) {
	once.Do(func() {
		cfg := config.Load()
		serverlessApp, serverlessErr = app.New(cfg)
	})
	return serverlessApp, serverlessErr
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fiberApp, err := getApp()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	adaptor.FiberApp(fiberApp).ServeHTTP(w, r)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
