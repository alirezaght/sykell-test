package url

import (
	"context"
	"sykell-backend/internal/db"
)

// RemoveURL deletes a URL by its ID for the specified user
func (s *Service) RemoveURL(ctx context.Context, userID string, urlID string) error {
	queries := db.New(s.db)
	err := queries.DeleteURLByIdAndUserId(ctx, db.DeleteURLByIdAndUserIdParams{
		ID:     urlID,
		UserID: userID,
	})
	return err
}