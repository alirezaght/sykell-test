package url

import (
	"context"
)

// FindUrls retrieves URLs based on the provided dashboard filters
func (s *Service) FindUrls(ctx context.Context, userID string, filters DashboardFilters) (PaginatedUrls, error) {
	
	// Map frontend sort column names to backend column names
	sortBy := mapSortColumn(filters.SortBy)
	sortDir := filters.SortOrder
	if sortDir == "" {
		sortDir = "desc"
	}
	
	crawlResults, err := s.repo.GetUrlsWithLatestCrawlsFiltered(ctx, userID, max(filters.Limit, 1), max(filters.Limit * (filters.Page - 1), 0), sortBy, sortDir, filters.Query)
	
	if err != nil {
		return PaginatedUrls{}, err
	}
	totalCount, err := s.repo.CountURLsByFilter(ctx, userID, filters.Query)
	if err != nil {
		return PaginatedUrls{}, err
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