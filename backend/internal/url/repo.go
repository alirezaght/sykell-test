package url

import (
	"context"
	"database/sql"
	"sykell-backend/internal/config"
	"sykell-backend/internal/db"
)

// Repo defines the interface for URL repository operations
type Repo interface {
	RemoveURL(ctx context.Context, userID string, urlID string) error
	CreateURL(ctx context.Context, userID string, normalizedURL string, domain string) error
	CountURLsByUserID(ctx context.Context, userID string) (int64, error)
	GetUrlsWithLatestCrawlsFiltered(ctx context.Context, userID string, limit int32, offset int32, sortBy string, sortOrder string, filter string) ([]CrawlResult, error)
}

type urlRepo struct {
	sqlDB *sql.DB
}

// NewRepo creates a new instance of the URL repository
func NewRepo(db *sql.DB) Repo {
	return &urlRepo{
		sqlDB: db,
	}
}

// RemoveURL deletes a URL by its ID for the specified user
func (r *urlRepo) RemoveURL(ctx context.Context, userID string, urlID string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	err := queries.DeleteURLByIdAndUserId(ctx, db.DeleteURLByIdAndUserIdParams{
		ID:     urlID,
		UserID: userID,
	})
	return err
}

// CreateURL creates a new URL entry for the specified user
func (r *urlRepo) CreateURL(ctx context.Context, userID string, normalizedURL string, domain string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	_, err := queries.CreateUrl(ctx, db.CreateUrlParams{
		UserID: 	userID,
		NormalizedUrl: normalizedURL,
		Domain: domain,
	})
	
	return err
}

// CountURLsByUserID counts the number of URLs for a given user
func (r *urlRepo) CountURLsByUserID(ctx context.Context, userID string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	count, err := queries.CountUrlsByUser(ctx, userID)
	return count, err
}

// GetUrlsWithLatestCrawlsFiltered retrieves URLs with their latest crawl results based on filters
func (r *urlRepo) GetUrlsWithLatestCrawlsFiltered(ctx context.Context, userID string, limit int32, offset int32, sortBy string, sortOrder string, filter string) ([]CrawlResult, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	result, err := queries.GetUrlsWithLatestCrawlsFiltered(ctx, db.GetUrlsWithLatestCrawlsFilteredParams{
		UserID:      userID,
		QueryFilter: filter,
		SortBy:      sortBy,
		SortDir:     sortOrder,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, err
	}
	crawlResults := make([]CrawlResult, len(result))
	for i, row := range result {		
		crawlResults[i] = convertDbRowToCrawlResult(row)
	}
	return crawlResults, nil
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