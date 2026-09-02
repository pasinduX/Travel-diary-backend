package handler

import (
	"strings"
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/middleware"
	"travel-diary-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func PricingListHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	_ = cfg
	svc := service.NewPricingService(dao.NewPricingDAO(db))
	return func(c *fiber.Ctx) error {
		plans, err := svc.List(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not load pricing plans"})
		}
		return c.JSON(fiber.Map{"data": plans})
	}
}

func PricingGetHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	_ = cfg
	svc := service.NewPricingService(dao.NewPricingDAO(db))
	return func(c *fiber.Ctx) error {
		plan, err := svc.Get(c.Context(), c.Params("slug"))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "pricing plan not found"})
		}
		return c.JSON(fiber.Map{"data": plan})
	}
}

func CurrentPlanHandler(db *mongo.Database) fiber.Handler {
	svc := service.NewPricingService(dao.NewPricingDAO(db), dao.NewUserDAO(db))
	return func(c *fiber.Ctx) error {
		userID, ok := middleware.Auth0UserID(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		plan, err := svc.GetForUser(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "could not load current pricing plan"})
		}
		return c.JSON(fiber.Map{"data": plan})
	}
}

func ChangePlanHandler(db *mongo.Database) fiber.Handler {
	svc := service.NewPricingService(dao.NewPricingDAO(db), dao.NewUserDAO(db))
	return func(c *fiber.Ctx) error {
		userID, ok := middleware.Auth0UserID(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}
		var req dto.PricingPlanRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		plan, err := svc.ChangeForUser(c.Context(), userID, strings.TrimSpace(req.Slug))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "could not change pricing plan"})
		}
		return c.JSON(fiber.Map{"data": plan})
	}
}
