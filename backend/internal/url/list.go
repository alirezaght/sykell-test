package url

import (
	"context"
	"sykell-backend/internal/db"
)

// FindUrls retrieves URLs based on the provided dashboard filters
func (s *Service) FindUrls(ctx context.Context, userID string, filters DashboardFilters) (PaginatedUrls, error) {
	queries := db.New(s.db)
	result, err := queries.GetUrlsWithLatestCrawlsFiltered(ctx, db.GetUrlsWithLatestCrawlsFilteredParams{
		UserID:     userID,
		QueryFilter: filters.Query,		
		Limit:      max(filters.Limit, 1),
		Offset:     max(filters.Limit * (filters.Page - 1), 0),
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

	return PaginatedUrls{
		Total: totalCount,
		Urls:  result,
		Page:  filters.Page,
		Limit: filters.Limit,
	}, nil
}