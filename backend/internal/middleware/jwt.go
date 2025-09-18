package middleware

import (
	"net/http"
	"strings"
	"github.com/labstack/echo/v4"
	"sykell-backend/internal/utils"
)

// JWTMiddleware creates a middleware function for JWT authentication
func JWTMiddleware(jwtSecret []byte) echo.MiddlewareFunc {
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			// Check if it starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			// Extract the token
			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing token",
				})
			}

			// Validate the token
			claims, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			// Store user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

