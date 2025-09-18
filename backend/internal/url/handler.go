package url

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Handler handles HTTP requests related to URLs
type Handler struct {
	urlService *Service
}


// NewHandler creates a new UrlHandler
func NewHandler(urlService *Service) *Handler {
	return &Handler{
		urlService: urlService,
	}
}


// AddURL handles adding a new URL
func (h *Handler) AddURL(c echo.Context) error {
	userID := c.Get("user_id")
	var req AddRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	ctx := c.Request().Context()

	err := h.urlService.AddURL(ctx, userID.(string), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve user profile",
		})
	}


	return c.NoContent(http.StatusOK)
}

// RemoveURL handles removing a URL
func (h *Handler) RemoveURL(c echo.Context) error {
	userID := c.Get("user_id")
	urlID := c.Param("id")
	if urlID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing URL ID",
		})
	}

	ctx := c.Request().Context()

	err := h.urlService.RemoveURL(ctx, userID.(string), urlID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to remove URL",
		})
	}

	return c.NoContent(http.StatusOK)
}


// ListURLs handles listing URLs with pagination
func (h *Handler) ListURLs(c echo.Context) error {
	userID := c.Get("user_id")
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	sortBy := c.QueryParam("sort_by")
	order := c.QueryParam("order")
	query := c.QueryParam("query")
	ctx := c.Request().Context()

	limitInt, _ := strconv.Atoi(limit)
	pageInt, _ := strconv.Atoi(page)

	filters := DashboardFilters{
		Query:    query,
		SortBy:  sortBy,
		SortOrder: order,
		Limit:    int32(limitInt),
		Page:    int32(pageInt),
	}

	result, err := h.urlService.FindUrls(ctx, userID.(string), filters)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to list URLs",
		})
	}
	return c.JSON(http.StatusOK, result)
}

	