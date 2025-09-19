package middleware

import (
	"time"

	"sykell-backend/internal/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ZapLogger returns a middleware that logs HTTP requests using Zap
func ZapLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			err := next(c)
			
			req := c.Request()
			res := c.Response()
			
			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("remote_ip", c.RealIP()),
				zap.String("user_agent", req.UserAgent()),
				zap.Int("status", res.Status),
				zap.Int64("bytes_out", res.Size),
				zap.Duration("latency", time.Since(start)),
			}
			
			// Add error to log if present
			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Error("Request completed with error", fields...)
			} else {
				// Log at different levels based on status code
				if res.Status >= 500 {
					logger.Error("Request completed", fields...)
				} else if res.Status >= 400 {
					logger.Warn("Request completed", fields...)
				} else {
					logger.Info("Request completed", fields...)
				}
			}
			
			return err
		}
	}
}