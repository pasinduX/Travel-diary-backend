package router

import (
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/handler"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(app *fiber.App, cfg config.Config, db *mongo.Database) {
	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/health", handler.HealthCheck(cfg))
	auth := v1.Group("/auth")
	auth.Post("/register", handler.RegisterHandler(cfg, db))
	auth.Post("/login", handler.LoginHandler(cfg, db))
	auth.Post("/refresh", handler.RefreshHandler(cfg, db))
	auth.Get("/google", handler.GoogleLoginHandler(cfg, db))
	auth.Get("/google/callback", handler.GoogleCallbackHandler(cfg, db))

	trips := v1.Group("/trips")
	trips.Post("/", handler.TripCreateHandler(cfg, db))
	trips.Get("/", handler.TripListHandler(cfg, db))
	trips.Get("/:id", handler.TripGetHandler(cfg, db))
	trips.Put("/:id", handler.TripUpdateHandler(cfg, db))
	trips.Delete("/:id", handler.TripDeleteHandler(cfg, db))
	trips.Post("/:id/images", handler.TripImageUploadHandler(cfg, db))
	trips.Get("/:id/images", handler.TripImageListHandler(cfg, db))
}
