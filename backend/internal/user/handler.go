package user

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{
		userService: userService,		
	}
}


// Login handles user authentication
func (h *UserHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	ctx := c.Request().Context()
	response, err := h.userService.Login(ctx, req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
	}

	// Set HTTP-only cookie with the JWT token
	cookie := &http.Cookie{
		Name:     "token",
		Value:    response.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		// Don't set SameSite for development - let browser decide
		Expires:  time.Unix(response.ExpiresAt, 0),
	}
	c.SetCookie(cookie)
	log.Printf("Login: Set cookie for user, token length: %d", len(response.Token))

	return c.JSON(http.StatusOK, response)
}


// Register handles new user registration
func (h *UserHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	ctx := c.Request().Context()
	err := h.userService.Register(ctx, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

// GetProfile retrieves the profile of the authenticated user
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID := c.Get("user_id")	

	ctx := c.Request().Context()

	user, err := h.userService.GetProfile(ctx, userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user profile",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// Logout handles user logout by clearing the token cookie
func (h *UserHandler) Logout(c echo.Context) error {
	// Clear the token cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0), // Set to past date to delete cookie
		MaxAge:   -1,
	}
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Logged out successfully",
	})
}
