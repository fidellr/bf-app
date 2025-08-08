package middleware

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type contextKey string

const (
	TraceIDKey contextKey = "trace_id"
)

func Tracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			traceID := uuid.New().String()

			ctx := context.WithValue(c.Request().Context(), TraceIDKey, traceID)
			c.SetRequest(c.Request().WithContext(ctx))

			c.Response().Header().Set("X-Trace-ID", traceID)

			zap.L().Info("request started",
				zap.String("trace_id", traceID),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Path()),
			)

			err := next(c)

			zap.L().Info("request completed",
				zap.String("trace_id", traceID),
				zap.Int("status", c.Response().Status),
				zap.Duration("duration", time.Since(start)),
			)

			return err
		}
	}
}
