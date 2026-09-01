package handler

import (
	"strings"
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/integrations"
	"travel-diary-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func TripCreateHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewTripService(dao.NewTripDAO(db))
	return func(c *fiber.Ctx) error {
		userID, err := userIDFromAccessToken(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		var req dto.TripCreateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		resp, err := svc.Create(c.Context(), userID, req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(resp)
	}
}

func TripListHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewTripService(dao.NewTripDAO(db))
	return func(c *fiber.Ctx) error {
		userID, err := userIDFromAccessToken(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		resp, err := svc.List(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	}
}

func TripGetHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewTripService(dao.NewTripDAO(db))
	return func(c *fiber.Ctx) error {
		userID, err := userIDFromAccessToken(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		resp, err := svc.Get(c.Context(), userID, c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	}
}

func TripUpdateHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewTripService(dao.NewTripDAO(db))
	return func(c *fiber.Ctx) error {
		userID, err := userIDFromAccessToken(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		var req dto.TripUpdateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		resp, err := svc.Update(c.Context(), userID, c.Params("id"), req)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	}
}

func TripDeleteHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewTripService(dao.NewTripDAO(db))
	return func(c *fiber.Ctx) error {
		userID, err := userIDFromAccessToken(cfg, c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		if err := svc.Delete(c.Context(), userID, c.Params("id")); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

func userIDFromAccessToken(cfg config.Config, c *fiber.Ctx) (string, error) {
	token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
	if token == "" {
		return "", fiber.ErrUnauthorized
	}
	claims, err := integrations.ParseAccessToken(cfg.JWTSecret, token)
	if err != nil {
		return "", fiber.ErrUnauthorized
	}
	userID, _ := claims["sub"].(string)
	if userID == "" {
		return "", fiber.ErrUnauthorized
	}
	return userID, nil
}
