package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/middleware"
)

func TripImageAnalysisHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := middleware.Auth0UserID(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		analysis, err := dao.NewTripImageAnalysisDAO(db).FindByImageID(
			c.Context(), userID, c.Params("id"), c.Params("imageId"),
		)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "image analysis not found"})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not load image analysis"})
		}

		return c.JSON(analysis)
	}
}

func TripAnalysisStatusHandler(db *mongo.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := middleware.Auth0UserID(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		counts, total, err := dao.NewTripImageDAO(db).CountByStatus(c.Context(), userID, c.Params("id"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "could not load analysis status"})
		}
		analyzed := counts["ANALYZED"]
		failed := counts["FAILED"]
		images, err := dao.NewTripImageDAO(db).ListByTripID(c.Context(), userID, c.Params("id"))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "could not load analysis failures"})
		}
		failures := make([]fiber.Map, 0)
		for _, image := range images {
			if image.AnalysisStatus == "FAILED" {
				failures = append(failures, fiber.Map{"imageId": image.ID, "fileName": image.FileName, "error": image.AnalysisError})
			}
		}
		return c.JSON(fiber.Map{"tripId": c.Params("id"), "total": total, "uploaded": counts["UPLOADED"], "queued": counts["QUEUED"], "processing": counts["PROCESSING"], "analyzed": analyzed, "failed": failed, "failures": failures, "percentage": func() int {
			if total == 0 {
				return 0
			}
			return int(analyzed * 100 / total)
		}(), "readyToGenerate": total > 0 && analyzed+failed == total})
	}
}
