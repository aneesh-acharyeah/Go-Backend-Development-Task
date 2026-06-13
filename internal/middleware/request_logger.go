package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const RequestIDHeader = "X-Request-Id"

func RequestLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		requestID := c.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		c.Set(RequestIDHeader, requestID)
		c.Locals("request_id", requestID)

		err := c.Next()
		duration := time.Since(start)
		status := c.Response().StatusCode()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", status),
			zap.Int64("duration_ms", duration.Milliseconds()),
		}

		if err != nil {
			logger.Error("request failed", append(fields, zap.Error(err))...)
			return err
		}

		if status >= fiber.StatusInternalServerError {
			logger.Error("request completed", fields...)
			return nil
		}

		logger.Info("request completed", fields...)
		return nil
	}
}
