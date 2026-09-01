package handler

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
)

func TripAnalysisStatusHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := userIDFromAccessToken(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		counts, total, err := dao.NewTripImageDAO(db).CountByStatus(c.Context(), userID, c.Params("id"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "could not load analysis status"})
		}
		analyzed := counts["ANALYZED"]
		failed := counts["FAILED"]
		return c.JSON(fiber.Map{"tripId": c.Params("id"), "total": total, "uploaded": counts["UPLOADED"], "queued": counts["QUEUED"], "processing": counts["PROCESSING"], "analyzed": analyzed, "failed": failed, "percentage": func() int {
			if total == 0 {
				return 0
			}
			return int(analyzed * 100 / total)
		}(), "readyToGenerate": total > 0 && analyzed+failed == total})
	}
}
