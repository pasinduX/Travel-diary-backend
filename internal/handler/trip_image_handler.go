package handler

import (
	"context"

	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/integrations"
	"travel-diary-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func TripImageUploadHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc, err := newTripImageService(cfg, db)
	if err != nil {
		panic(err)
	}

	return func(c *fiber.Ctx) error {
		userID, err := userIDFromAccessToken(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		tripID := c.Params("id")

		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "multipart form required"})
		}
		files := form.File["images"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no images provided"})
		}

		resp, err := svc.UploadMany(c.Context(), userID, tripID, files)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(dto.TripImageUploadResponse{Uploaded: resp})
	}
}

func TripImageListHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc, err := newTripImageService(cfg, db)
	if err != nil {
		panic(err)
	}

	return func(c *fiber.Ctx) error {
		userID, err := userIDFromAccessToken(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		resp, err := svc.List(c.Context(), userID, c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(dto.TripImageUploadResponse{Uploaded: resp})
	}
}

func newTripImageService(cfg config.Config, db *mongo.Database) (*service.TripImageService, error) {
	s3Client, err := integrations.NewS3Client(context.Background(), cfg.AWSRegion, cfg.AWSAccessKeyID, cfg.AWSSecretAccessKey, cfg.AWSBucket)
	if err != nil {
		return nil, err
	}
	return service.NewTripImageService(dao.NewTripDAO(db), dao.NewTripImageDAO(db), s3Client), nil
}
