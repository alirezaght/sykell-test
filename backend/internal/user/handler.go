package user

import (
	"net/http"
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
	_, err := h.userService.Register(ctx, req)
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
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Unauthorized",
		})
	}

	ctx := c.Request().Context()

	user, err := h.userService.GetProfile(ctx, userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user profile",
		})
	}


	return c.JSON(http.StatusOK, user)
}
