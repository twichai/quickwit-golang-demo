package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RequestLogger(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		latency := time.Since(start)

		log.Info("http_request",
			zap.String("title", "HTTP Request"),
			zap.String("body", c.Method()+" "+c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", latency),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}
