package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Recover() func(*fiber.Ctx) error {
	return recover.New()
}
