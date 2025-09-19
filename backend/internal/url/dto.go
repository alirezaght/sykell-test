package url

import (
	"time"
)

// CrawlResult represents a URL with its latest crawl information for frontend consumption
type CrawlResult struct {
	UrlID                  string    `json:"url_id"`
	NormalizedUrl          string    `json:"normalized_url"`
	Domain                 string    `json:"domain"`
	UrlCreatedAt           *time.Time `json:"url_created_at"`
	CrawlID                *string   `json:"crawl_id"`
	Status                 *string   `json:"status"`
	QueuedAt               *time.Time `json:"queued_at"`
	StartedAt              *time.Time `json:"started_at"`
	FinishedAt             *time.Time `json:"finished_at"`
	HtmlVersion            *string   `json:"html_version"`
	PageTitle              *string   `json:"page_title"`
	H1Count                *int32    `json:"h1_count"`
	H2Count                *int32    `json:"h2_count"`
	H3Count                *int32    `json:"h3_count"`
	H4Count                *int32    `json:"h4_count"`
	H5Count                *int32    `json:"h5_count"`
	H6Count                *int32    `json:"h6_count"`
	InternalLinksCount     *int32    `json:"internal_links_count"`
	ExternalLinksCount     *int32    `json:"external_links_count"`
	InaccessibleLinksCount *int32    `json:"inaccessible_links_count"`
	HasLoginForm           *bool     `json:"has_login_form"`
	ErrorMessage           *string   `json:"error_message"`
	CrawlCreatedAt         *time.Time `json:"crawl_created_at"`
	CrawlUpdatedAt         *time.Time `json:"crawl_updated_at"`
}

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
	Total int64         `json:"total_count"`
	Urls  []CrawlResult `json:"urls"`
	Page  int32         `json:"page"`
	Limit int32         `json:"limit"`
}

// AddRequest represents the request payload for adding a new URL
type AddRequest struct {
	URL string `json:"url" validate:"required,url"`
}


