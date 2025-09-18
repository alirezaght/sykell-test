package url

import (
	"sykell-backend/internal/db"
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

// convertDbRowToCrawlResult converts a database row to a CrawlResult
func convertDbRowToCrawlResult(row db.GetUrlsWithLatestCrawlsFilteredRow) CrawlResult {
	result := CrawlResult{
		UrlID:         row.UrlID,
		NormalizedUrl: row.NormalizedUrl,
		Domain:        row.Domain,
	}

	// Convert nullable time
	if row.UrlCreatedAt.Valid {
		result.UrlCreatedAt = &row.UrlCreatedAt.Time
	}

	// Convert nullable strings
	if row.CrawlID.Valid {
		result.CrawlID = &row.CrawlID.String
	}
	if row.Status.Valid {
		statusStr := string(row.Status.CrawlsStatus)
		result.Status = &statusStr
	}
	if row.HtmlVersion.Valid {
		result.HtmlVersion = &row.HtmlVersion.String
	}
	if row.PageTitle.Valid {
		result.PageTitle = &row.PageTitle.String
	}
	if row.ErrorMessage.Valid {
		result.ErrorMessage = &row.ErrorMessage.String
	}

	// Convert nullable times
	if row.QueuedAt.Valid {
		result.QueuedAt = &row.QueuedAt.Time
	}
	if row.StartedAt.Valid {
		result.StartedAt = &row.StartedAt.Time
	}
	if row.FinishedAt.Valid {
		result.FinishedAt = &row.FinishedAt.Time
	}
	if row.CrawlCreatedAt.Valid {
		result.CrawlCreatedAt = &row.CrawlCreatedAt.Time
	}
	if row.CrawlUpdatedAt.Valid {
		result.CrawlUpdatedAt = &row.CrawlUpdatedAt.Time
	}

	// Convert nullable integers
	if row.H1Count.Valid {
		result.H1Count = &row.H1Count.Int32
	}
	if row.H2Count.Valid {
		result.H2Count = &row.H2Count.Int32
	}
	if row.H3Count.Valid {
		result.H3Count = &row.H3Count.Int32
	}
	if row.H4Count.Valid {
		result.H4Count = &row.H4Count.Int32
	}
	if row.H5Count.Valid {
		result.H5Count = &row.H5Count.Int32
	}
	if row.H6Count.Valid {
		result.H6Count = &row.H6Count.Int32
	}
	if row.InternalLinksCount.Valid {
		result.InternalLinksCount = &row.InternalLinksCount.Int32
	}
	if row.ExternalLinksCount.Valid {
		result.ExternalLinksCount = &row.ExternalLinksCount.Int32
	}
	if row.InaccessibleLinksCount.Valid {
		result.InaccessibleLinksCount = &row.InaccessibleLinksCount.Int32
	}

	// Convert nullable bool
	if row.HasLoginForm.Valid {
		result.HasLoginForm = &row.HasLoginForm.Bool
	}

	return result
}
