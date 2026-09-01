package handler

import (
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/service"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(cfg config.Config) fiber.Handler {
	svc := service.NewHealthService(cfg)

	return func(c *fiber.Ctx) error {
		resp := svc.Check(c.Context())
		return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
			Message: "ok",
			Data:    resp,
		})
	}
}
