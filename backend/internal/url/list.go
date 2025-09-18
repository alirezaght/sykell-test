package url

import (
	"context"
	"sykell-backend/internal/db"
)

// FindUrls retrieves URLs based on the provided dashboard filters
func (s *Service) FindUrls(ctx context.Context, userID string, filters DashboardFilters) (PaginatedUrls, error) {
	queries := db.New(s.db)
	
	// Map frontend sort column names to backend column names
	sortBy := mapSortColumn(filters.SortBy)
	sortDir := filters.SortOrder
	if sortDir == "" {
		sortDir = "desc"
	}
	
	result, err := queries.GetUrlsWithLatestCrawlsFiltered(ctx, db.GetUrlsWithLatestCrawlsFilteredParams{
		UserID:      userID,
		QueryFilter: filters.Query,
		SortBy:      sortBy,
		SortDir:     sortDir,
		Limit:       max(filters.Limit, 1),
		Offset:      max(filters.Limit * (filters.Page - 1), 0),
	})
	if err != nil {
		return PaginatedUrls{}, err
	}
	totalCount, err := queries.CountUrlsWithFilter(ctx, db.CountUrlsWithFilterParams{
		UserID:     userID,
		QueryFilter: filters.Query,
	})
	if err != nil {
		return PaginatedUrls{}, err
	}

	// Convert database rows to frontend-friendly format
	crawlResults := make([]CrawlResult, len(result))
	for i, row := range result {
		crawlResults[i] = convertDbRowToCrawlResult(row)
	}

	return PaginatedUrls{
		Total: totalCount,
		Urls:  crawlResults,
		Page:  filters.Page,
		Limit: filters.Limit,
	}, nil
}

// mapSortColumn maps frontend sort column names to backend column names
func mapSortColumn(frontendColumn string) string {
	columnMap := map[string]string{
		"url":                 "normalized_url",
		"domain":              "domain",
		"title":               "page_title",
		"status":              "status",
		"html_version":        "html_version",
		"internal_links":      "internal_links_count",
		"external_links":      "external_links_count",
		"inaccessible_links":  "inaccessible_links_count",
		"h1_count":            "h1_count",
		"h2_count":            "h2_count",
		"h3_count":            "h3_count",
		"created_at":          "url_created_at",
		"finished_at":         "finished_at",
	}
	
	if backendColumn, exists := columnMap[frontendColumn]; exists {
		return backendColumn
	}
	
	// Default sort column
	return "url_created_at"
}