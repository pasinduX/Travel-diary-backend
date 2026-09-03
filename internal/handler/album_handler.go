package handler

import (
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/middleware"
	"travel-diary-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateAlbumHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewAlbumService(dao.NewTripDAO(db), dao.NewTripImageDAO(db), dao.NewTripImageAnalysisDAO(db), dao.NewAlbumPlanDAO(db), cfg)
	return func(c *fiber.Ctx) error {
		userID, ok := middleware.Auth0UserID(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		plan, err := svc.Generate(c.Context(), userID, c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(plan.Plan)
	}
}

func AlbumGetHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewAlbumService(dao.NewTripDAO(db), dao.NewTripImageDAO(db), dao.NewTripImageAnalysisDAO(db), dao.NewAlbumPlanDAO(db), cfg)
	return func(c *fiber.Ctx) error {
		userID, ok := middleware.Auth0UserID(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		plan, err := svc.Get(c.Context(), userID, c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "album not generated"})
		}
		return c.JSON(plan.Plan)
	}
}
