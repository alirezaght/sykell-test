package user

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	// db *sql.DB - will be added when we connect to database
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateUserRequest struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	// Mock data for now
	users := []User{
		{ID: 1, Email: "john@example.com", Name: "John Doe"},
		{ID: 2, Email: "jane@example.com", Name: "Jane Smith"},
	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Mock data for now
	user := User{
		ID:    id,
		Email: "user@example.com",
		Name:  "Sample User",
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Mock response for now
	user := User{
		ID:    3,
		Email: req.Email,
		Name:  req.Name,
	}

	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	var req CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Mock response for now
	user := User{
		ID:    id,
		Email: req.Email,
		Name:  req.Name,
	}

	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID",
		})
	}

	// Mock response for now
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
		"id":      id,
	})
}