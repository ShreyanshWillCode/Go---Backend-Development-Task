package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/shreyxnsh/anyx-user-api/internal/logger"
	"go.uber.org/zap"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start)

		requestID, _ := c.Locals("requestID").(string)

		logger.Info("request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
		)

		return err
	}
}
