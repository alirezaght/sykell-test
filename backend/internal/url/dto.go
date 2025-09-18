package url

import (
	"sykell-backend/internal/db"
)


// DashboardFilters represents the filtering options for the dashboard
type DashboardFilters struct {
	Query   string `json:"query"`	
	SortBy      string `json:"sort_by"`
	SortOrder   string `json:"sort_order"` // "asc" or "desc"
	Limit       int32  `json:"limit"`
	Page      int32  `json:"offset"`
}

// PaginatedUrls represents a paginated list of URLs with metadata
type PaginatedUrls struct {
	Total int64      `json:"total_count"`
	Urls       []db.GetUrlsWithLatestCrawlsFilteredRow `json:"urls"`
	Page 	 int32      `json:"page"`
	Limit   int32      `json:"limit"`
}

// AddRequest represents the request payload for adding a new URL
type AddRequest struct {
	URL string `json:"url" validate:"required,url"`
}
