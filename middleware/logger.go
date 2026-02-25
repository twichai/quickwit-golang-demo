package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid" // แนะนำให้ลงเพิ่ม: go get github.com/google/uuid
	"go.uber.org/zap"
)

func RequestLogger(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// 1. Manage request codes (create a new code if one doesn't exist).
		rid := c.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}
		c.Set("X-Request-ID", rid)

		// 2. Go and work on other parts of the app.
		err := c.Next()

		// 3. Gather information after the work is completed.
		latency := time.Since(start)
		statusCode := c.Response().StatusCode()

		// 4. Prepare the data to match Quickwit YAML.
		fields := []zap.Field{
			zap.String("request_id", rid),
			zap.Int("status", statusCode),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int64("latency", latency.Nanoseconds()), // Send as an i64 number according to YAML.
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		}

		// 5. Decide whether to log as Info or Error.
		msg := "http_request"
		if err != nil {
			fields = append(fields, zap.Error(err))
			log.Error(msg, fields...)
		} else {
			log.Info(msg, fields...)
		}

		return err
	}
}
