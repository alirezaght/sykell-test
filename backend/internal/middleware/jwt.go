package middleware

import (
	"net/http"
	"strings"
	"sykell-backend/internal/logger"
	"sykell-backend/internal/utils"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// JWTMiddleware creates a middleware function for JWT authentication
func JWTMiddleware(jwtSecret []byte, fromCookie bool) echo.MiddlewareFunc {
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var token string
			
			// Debug logging
			logger.Debug("JWT middleware processing request", 
				zap.String("path", c.Request().URL.Path), 
				zap.Bool("from_cookie", fromCookie))
			
			if fromCookie {
				// Only check cookie
				cookie, err := c.Cookie("token")
				if err == nil && cookie.Value != "" {
					token = cookie.Value
					logger.Debug("JWT middleware using cookie token")
				}
				
				if token == "" {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "Missing token cookie",
					})
				}
			} else {
				// Only check Authorization header
				authHeader := c.Request().Header.Get("Authorization")
				// NEVER log the actual auth header as it contains the token
				logger.Debug("JWT middleware checking Authorization header", zap.Bool("has_auth_header", authHeader != ""))
				
				if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
					logger.Debug("JWT middleware using Bearer token")
				}
				
				if token == "" {
					return c.JSON(http.StatusUnauthorized, map[string]string{
						"error": "Missing or invalid Authorization header",
					})
				}
			}

			// Validate the token
			claims, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				logger.Warn("JWT middleware token validation failed", zap.Error(err))
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			logger.Debug("JWT middleware authentication successful", zap.String("user_id", claims.UserID))

			// Store user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

