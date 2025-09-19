package crawl

import (
	"context"
	"database/sql"
	"sykell-backend/internal/config"
	"sykell-backend/internal/db"
)


type Repo interface {
	GetCrawlIDByWorkflowID(ctx context.Context, workflowID string) (string, error)
	QueueCrawl(ctx context.Context, urlID string, workflowID string) error
	CountOfActiveCrawlForUrlId(ctx context.Context, urlID string) (int64, error)
	GetUrlByIdAndUserId(ctx context.Context, urlID string, userID string) (*URLResponse, error)
	UpdateCrawlResult(ctx context.Context, crawlID string, htmlVersion string, pageTitle string, h1Count int32, h2Count int32, h3Count int32, h4Count int32, h5Count int32, h6Count int32, internalLinksCount int32, externalLinksCount int32, inaccessableLinksCount int32, hasLoginForm bool, status string) error
	CreateInaccessibleLink(ctx context.Context, crawlID string, href string, absoluteURL string, isInternal bool, statusCode int, anchorText string) error
	SetCrawlError(ctx context.Context, crawlID string, errorMessage string) error
	SetCrawlRunning(ctx context.Context, crawlID string) error
	SetCrawlStopped(ctx context.Context, crawlID string) error
	GetActiveCrawlsForUrlId(ctx context.Context, urlID string) ([]CrawlResponse, error) 
}

type crawlRepo struct {
	sqlDB *sql.DB
}

func NewRepo(db *sql.DB) Repo {
	return &crawlRepo{
		sqlDB: db,
	}
}


func (r *crawlRepo) GetCrawlIDByWorkflowID(ctx context.Context, workflowID string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	crawl, err := queries.GetCrawlByWorkflowID(ctx, workflowID)
	if err != nil {
		return "", err
	}
	return crawl.ID, nil
}

func (r *crawlRepo) QueueCrawl(ctx context.Context, urlID string, workflowID string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	_, err := queries.QueueCrawl(ctx, db.QueueCrawlParams{
		UrlID: urlID,
		WorkflowID: workflowID,
	})
	return err
}

func (r *crawlRepo) CountOfActiveCrawlForUrlId(ctx context.Context, urlID string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	count, err := queries.CountOfActiveCrawlForUrlId(ctx, urlID)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *crawlRepo) GetUrlByIdAndUserId(ctx context.Context, urlID string, userID string) (*URLResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	url, err := queries.GetUrlByIdAndUserId(ctx, db.GetUrlByIdAndUserIdParams{
		ID:   urlID,
		UserID: userID,
	})
	if err != nil {		
		return &URLResponse{}, err
	}
	return &URLResponse{
		ID: url.ID,
		Domain: url.Domain,
		NormalizedUrl: url.NormalizedUrl,
	}, nil
}

func (r *crawlRepo) UpdateCrawlResult(ctx context.Context, crawlID string, htmlVersion string, pageTitle string, h1Count int32, h2Count int32, h3Count int32, h4Count int32, h5Count int32, h6Count int32, internalLinksCount int32, externalLinksCount int32, inaccessableLinksCount int32, hasLoginForm bool, status string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	err := queries.UpdateCrawlResult(ctx, db.UpdateCrawlResultParams{
		ID:                    crawlID,
		HtmlVersion:           sql.NullString{String: htmlVersion, Valid: htmlVersion != ""},
		PageTitle:             sql.NullString{String: pageTitle, Valid: pageTitle != ""},
		H1Count:               sql.NullInt32{Int32: h1Count, Valid: true},
		H2Count:               sql.NullInt32{Int32: h2Count, Valid: true},
		H3Count:               sql.NullInt32{Int32: h3Count, Valid: true},
		H4Count:               sql.NullInt32{Int32: h4Count, Valid: true},
		H5Count:               sql.NullInt32{Int32: h5Count, Valid: true},
		H6Count:               sql.NullInt32{Int32: h6Count, Valid: true},
		InternalLinksCount:    sql.NullInt32{Int32: internalLinksCount, Valid: true},
		ExternalLinksCount:    sql.NullInt32{Int32: externalLinksCount, Valid: true},
		InaccessibleLinksCount: sql.NullInt32{Int32: inaccessableLinksCount, Valid: true},
		HasLoginForm:          hasLoginForm,
		Status:				db.CrawlsStatus(status),
	})
	return err
}

func (r *crawlRepo) CreateInaccessibleLink(ctx context.Context, crawlID string, href string, absoluteURL string, isInternal bool, statusCode int, anchorText string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	_, err := queries.CreateInaccessibleLink(ctx, db.CreateInaccessibleLinkParams{
		CrawlID:     crawlID,
		Href:        href,
		AbsoluteUrl: absoluteURL,
		IsInternal:  isInternal,
		StatusCode:  sql.NullInt32{Int32: int32(statusCode), Valid: true},
		AnchorText:  sql.NullString{String: anchorText, Valid: anchorText != ""},
	})
	return err
}

func (r *crawlRepo) SetCrawlError(ctx context.Context, crawlID string, errorMessage string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	err := queries.SetCrawlError(ctx, db.SetCrawlErrorParams{
		ID:           crawlID,
		ErrorMessage: sql.NullString{String: errorMessage, Valid: errorMessage != ""},
	})
	return err
}

func (r *crawlRepo) SetCrawlRunning(ctx context.Context, crawlID string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	err := queries.SetCrawlRunning(ctx, crawlID)
	return err
}

func (r *crawlRepo) SetCrawlStopped(ctx context.Context, crawlID string) error {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	err := queries.SetCrawlStopped(ctx, crawlID)
	return err
}

func (r *crawlRepo) GetActiveCrawlsForUrlId(ctx context.Context, urlID string) ([]CrawlResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultTimeout)
	defer cancel()
	queries := db.New(r.sqlDB)
	result, err := queries.GetActiveCrawlsForUrlId(ctx, urlID)
	if err != nil {
		return nil, err
	}
	crawls := make([]CrawlResponse, len(result))
	for i, row := range result {
		crawls[i] = CrawlResponse{
			ID: row.ID,
			WorkflowID: row.WorkflowID,
		}
	}
	return crawls, err
}