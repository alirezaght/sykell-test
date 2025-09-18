package url

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sykell-backend/internal/db"
)


func normalizeURL(raw string) (string, string, error) {
	// Parse
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}

	// Normalize scheme → lowercase
	scheme := strings.ToLower(u.Scheme)
	if scheme == "" {
		scheme = "http" // default if missing
	}

	// Normalize host → lowercase, strip default ports
	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if (scheme == "http" && port == "80") || (scheme == "https" && port == "443") {
		port = ""
	}

	// Rebuild authority (host[:port])
	domain := host
	if port != "" {
		domain = fmt.Sprintf("%s:%s", host, port)
	}

	// Normalize path → ensure it’s not empty
	path := u.EscapedPath()
	if path == "" {
		path = "/"
	}

	// Rebuild normalized URL (without fragment)
	normalized := scheme + "://" + domain + path
	if u.RawQuery != "" {
		normalized += "?" + u.RawQuery
	}

	return normalized, host, nil // domain = host only, no port
}

// AddURL adds a new URL for the specified user
func (s *Service) AddURL(ctx context.Context, userID string, request AddRequest) error {
	normalizeURL, domain, err := normalizeURL(request.URL)
	if err != nil {
		return err
	}
	queries := db.New(s.db)
	_, err = queries.CreateUrl(ctx, db.CreateUrlParams{
		UserID: userID,
		NormalizedUrl: normalizeURL,
		Domain: domain,		
	})
	return err
}