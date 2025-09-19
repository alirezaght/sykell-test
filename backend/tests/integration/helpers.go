package integration

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

// TestHelper provides common testing utilities for integration tests
type TestHelper struct {
	DB   *sql.DB
	Mock sqlmock.Sqlmock
	t    *testing.T
}

// NewTestHelper creates a new test helper with mock database
func NewTestHelper(t *testing.T) *TestHelper {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	
	return &TestHelper{
		DB:   mockDB,
		Mock: mock,
		t:    t,
	}
}

// Close cleans up the test helper resources
func (th *TestHelper) Close() {
	th.DB.Close()
}

// ExpectationsWereMet verifies all mock expectations were satisfied
func (th *TestHelper) ExpectationsWereMet() {
	require.NoError(th.t, th.Mock.ExpectationsWereMet())
}

// TestData contains common test data structures
type TestData struct {
	UserID     string
	URLID      string
	CrawlID    string
	WorkflowID string
	Domain     string
	URL        string
}

// NewTestData creates a new set of test data with realistic values
func NewTestData() *TestData {
	return &TestData{
		UserID:     "user_123e4567-e89b-12d3-a456-426614174000",
		URLID:      "url_123e4567-e89b-12d3-a456-426614174001",
		CrawlID:    "crawl_123e4567-e89b-12d3-a456-426614174002",
		WorkflowID: "workflow_123e4567-e89b-12d3-a456-426614174003",
		Domain:     "example.com",
		URL:        "https://example.com/test-page",
	}
}

// HTMLTestPage contains sample HTML for testing
const HTMLTestPage = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Integration Test Page</title>
    <meta name="description" content="A test page for integration testing">
</head>
<body>
    <header>
        <h1>Welcome to Test Page</h1>
        <nav>
            <ul>
                <li><a href="/home">Home</a></li>
                <li><a href="/about">About</a></li>
                <li><a href="/contact">Contact</a></li>
            </ul>
        </nav>
    </header>
    
    <main>
        <section>
            <h2>Main Content Section</h2>
            <p>This is a paragraph of content for testing purposes.</p>
            
            <h3>Subsection 1</h3>
            <p>Content for subsection 1.</p>
            
            <h3>Subsection 2</h3>
            <p>Content for subsection 2.</p>
            
            <h4>Sub-subsection</h4>
            <p>Nested content.</p>
        </section>
        
        <section>
            <h2>Links Section</h2>
            <p>Testing various types of links:</p>
            <ul>
                <li><a href="/internal-page">Internal link</a></li>
                <li><a href="https://external-site.com">External link</a></li>
                <li><a href="/broken-link">Potentially broken link</a></li>
                <li><a href="mailto:test@example.com">Email link</a></li>
                <li><a href="tel:+1234567890">Phone link</a></li>
            </ul>
        </section>
        
        <section>
            <h2>Forms Section</h2>
            
            <!-- Login form -->
            <form method="post" action="/login" class="login-form">
                <h3>Login Form</h3>
                <div>
                    <label for="email">Email:</label>
                    <input type="email" id="email" name="email" required>
                </div>
                <div>
                    <label for="password">Password:</label>
                    <input type="password" id="password" name="password" required>
                </div>
                <button type="submit">Login</button>
            </form>
            
            <!-- Contact form (not a login form) -->
            <form method="post" action="/contact" class="contact-form">
                <h3>Contact Form</h3>
                <div>
                    <label for="name">Name:</label>
                    <input type="text" id="name" name="name" required>
                </div>
                <div>
                    <label for="message">Message:</label>
                    <textarea id="message" name="message" required></textarea>
                </div>
                <button type="submit">Send Message</button>
            </form>
        </section>
    </main>
    
    <footer>
        <p>&copy; 2025 Test Site. All rights reserved.</p>
    </footer>
</body>
</html>`

// ExpectedCrawlResult contains the expected analysis results for HTMLTestPage
type ExpectedCrawlResult struct {
	HTMLVersion             string
	PageTitle               string
	H1Count                 int32
	H2Count                 int32
	H3Count                 int32
	H4Count                 int32
	H5Count                 int32
	H6Count                 int32
	InternalLinksCount      int32
	ExternalLinksCount      int32
	InaccessibleLinksCount  int32
	HasLoginForm            bool
}

// GetExpectedResult returns the expected crawl analysis for HTMLTestPage
func GetExpectedResult() *ExpectedCrawlResult {
	return &ExpectedCrawlResult{
		HTMLVersion:             "HTML5",
		PageTitle:               "Integration Test Page",
		H1Count:                 1, // "Welcome to Test Page"
		H2Count:                 3, // "Main Content Section", "Links Section", "Forms Section"
		H3Count:                 4, // "Subsection 1", "Subsection 2", "Login Form", "Contact Form"
		H4Count:                 1, // "Sub-subsection"
		H5Count:                 0,
		H6Count:                 0,
		InternalLinksCount:      4, // /home, /about, /contact, /internal-page, /broken-link
		ExternalLinksCount:      1, // https://external-site.com
		InaccessibleLinksCount:  0, // Would be determined during actual crawling
		HasLoginForm:            true, // Form with email/password fields
	}
}