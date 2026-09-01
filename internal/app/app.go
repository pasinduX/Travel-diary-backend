package app

import (
	"fmt"
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/integrations"
	"travel-diary-backend/internal/middleware"
	"travel-diary-backend/internal/router"

	"github.com/gofiber/fiber/v2"
	fibercors "github.com/gofiber/fiber/v2/middleware/cors"
)

func New(cfg config.Config) (*fiber.App, error) {
	// Fiber's default BodyLimit is 4 MB, which rejects even a single phone
	// photo on the trip-image multipart upload. Raise it app-wide.
	app := fiber.New(fiber.Config{
		BodyLimit: 64 * 1024 * 1024, // 64 MB
	})

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	client, err := integrations.NewMongoClient(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	app.Use(fibercors.New(fibercors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
	}))
	app.Use(middleware.RequestLogger())
	app.Use(middleware.Recover())

	router.Register(app, cfg, client.Database(cfg.MongoDatabase))

	return app, nil
}
