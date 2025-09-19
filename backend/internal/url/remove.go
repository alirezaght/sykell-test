package url

import (
	"context"	
)

// RemoveURL deletes a URL by its ID for the specified user
func (s *Service) RemoveURL(ctx context.Context, userID string, urlID string) error {	
	err := s.repo.RemoveURL(ctx, userID, urlID)
	return err
}