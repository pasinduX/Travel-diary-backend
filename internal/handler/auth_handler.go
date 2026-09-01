package handler

import (
	"net/url"
	"travel-diary-backend/internal/config"
	"travel-diary-backend/internal/dao"
	"travel-diary-backend/internal/dto"
	"travel-diary-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewAuthService(cfg, dao.NewUserDAO(db), dao.NewSessionDAO(db))

	return func(c *fiber.Ctx) error {
		var req dto.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		resp, err := svc.Register(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(resp)
	}
}

func LoginHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewAuthService(cfg, dao.NewUserDAO(db), dao.NewSessionDAO(db))

	return func(c *fiber.Ctx) error {
		var req dto.LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		resp, err := svc.Login(c.Context(), req)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(resp)
	}
}

func GoogleLoginHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewAuthService(cfg, dao.NewUserDAO(db), dao.NewSessionDAO(db))

	return func(c *fiber.Ctx) error {
		return c.Redirect(svc.GoogleLoginURL())
	}
}

func GoogleCallbackHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewAuthService(cfg, dao.NewUserDAO(db), dao.NewSessionDAO(db))

	return func(c *fiber.Ctx) error {
		code := c.Query("code")
		state := c.Query("state")
		if code == "" || state == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing code or state"})
		}

		resp, err := svc.GoogleCallback(c.Context(), code, state)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if cfg.FrontendRedirectURL != "" {
			return c.Redirect(cfg.FrontendRedirectURL + "?token=" + url.QueryEscape(resp.Token) + "&refreshToken=" + url.QueryEscape(resp.RefreshToken))
		}
		return c.JSON(resp)
	}
}

func RefreshHandler(cfg config.Config, db *mongo.Database) fiber.Handler {
	svc := service.NewAuthService(cfg, dao.NewUserDAO(db), dao.NewSessionDAO(db))

	return func(c *fiber.Ctx) error {
		var req dto.RefreshRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}
		if req.RefreshToken == "" {
			req.RefreshToken = c.Get("X-Refresh-Token")
		}
		if req.RefreshToken == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "refresh token required"})
		}

		resp, err := svc.Refresh(c.Context(), req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(resp)
	}
}
