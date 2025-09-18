package middleware

import (
	"log"
	"net/http"
	"strings"
	"sykell-backend/internal/utils"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware creates a middleware function for JWT authentication
func JWTMiddleware(jwtSecret []byte) echo.MiddlewareFunc {
	
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var token string
			
			// Debug logging
			log.Printf("JWT Middleware - Path: %s", c.Request().URL.Path)
			
			// First try to get from cookie (preferred method)
			cookie, err := c.Cookie("token")
			if err == nil && cookie.Value != "" {
				token = cookie.Value
				log.Printf("JWT Middleware - Using cookie token")
			} else {
				// Try Authorization header as fallback
				authHeader := c.Request().Header.Get("Authorization")
				log.Printf("JWT Middleware - Authorization header: %s", authHeader)
				
				if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimPrefix(authHeader, "Bearer ")
					log.Printf("JWT Middleware - Using Bearer token")
				} else {
					// Try query parameter as last resort
					queryToken := c.QueryParam("token")
					if queryToken != "" {
						token = queryToken
						log.Printf("JWT Middleware - Using query parameter token")
					}
				}
			}

			if token == "" {								
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header, token query parameter, or token cookie",
				})
			}

			// Validate the token
			claims, err := utils.ValidateJWT(token, jwtSecret)
			if err != nil {
				log.Printf("JWT Middleware - Token validation failed: %v", err)
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid token",
				})
			}

			log.Printf("JWT Middleware - Authentication successful for user: %s", claims.UserID)

			// Store user information in context
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

