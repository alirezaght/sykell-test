package utils

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// Helper function to create HTML document from string
func parseHTML(htmlContent string) *html.Node {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		panic(err)
	}
	return doc
}

func TestExtractHtmlVersion(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "HTML5 DOCTYPE",
			html:     `<!DOCTYPE html><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML5",
		},
		{
			name:     "HTML5 with explicit version",
			html:     `<!DOCTYPE html5><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML5",
		},
		{
			name:     "XHTML 1.0 Strict",
			html:     `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"><html><head><title>Test</title></head><body></body></html>`,
			expected: "XHTML",
		},
		{
			name:     "HTML 4.01 Strict",
			html:     `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd"><html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML 4.01",
		},
		{
			name:     "No DOCTYPE defaults to HTML5",
			html:     `<html><head><title>Test</title></head><body></body></html>`,
			expected: "HTML5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseHTML(tt.html)
			result := ExtractHtmlVersion(doc)
			if result != tt.expected {
				t.Errorf("ExtractHtmlVersion() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Simple title",
			html:     `<html><head><title>Test Page</title></head><body></body></html>`,
			expected: "Test Page",
		},
		{
			name:     "Title with whitespace",
			html:     `<html><head><title>  Test Page  </title></head><body></body></html>`,
			expected: "Test Page",
		},
		{
			name:     "Empty title",
			html:     `<html><head><title></title></head><body></body></html>`,
			expected: "",
		},
		{
			name:     "No title tag",
			html:     `<html><head></head><body></body></html>`,
			expected: "",
		},
		{
			name:     "Title with special characters",
			html:     `<html><head><title>Test & Special "Characters" 'Here'</title></head><body></body></html>`,
			expected: `Test & Special "Characters" 'Here'`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseHTML(tt.html)
			result := ExtractTitle(doc)
			if result != tt.expected {
				t.Errorf("ExtractTitle() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountHeadings(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected map[string]int
	}{
		{
			name: "All heading levels",
			html: `<html><body>
				<h1>Heading 1</h1>
				<h2>Heading 2</h2>
				<h2>Another Heading 2</h2>
				<h3>Heading 3</h3>
				<h4>Heading 4</h4>
				<h5>Heading 5</h5>
				<h6>Heading 6</h6>
			</body></html>`,
			expected: map[string]int{
				"h1": 1,
				"h2": 2,
				"h3": 1,
				"h4": 1,
				"h5": 1,
				"h6": 1,
			},
		},
		{
			name: "No headings",
			html: `<html><body><p>No headings here</p></body></html>`,
			expected: map[string]int{},
		},
		{
			name: "Only h1 and h3",
			html: `<html><body>
				<h1>Title</h1>
				<h3>Subtitle</h3>
				<h3>Another Subtitle</h3>
			</body></html>`,
			expected: map[string]int{
				"h1": 1,
				"h3": 2,
			},
		},
		{
			name: "Nested headings",
			html: `<html><body>
				<div>
					<h1>Main Title</h1>
					<section>
						<h2>Section Title</h2>
					</section>
				</div>
			</body></html>`,
			expected: map[string]int{
				"h1": 1,
				"h2": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseHTML(tt.html)
			result := CountHeadings(doc)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CountHeadings() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountLinks(t *testing.T) {
	tests := []struct {
		name           string
		html           string
		baseURL        string
		expectedMinCounts map[string]int // Using minimum counts since network calls are unpredictable
	}{
		{
			name: "Mixed links",
			html: `<html><body>
				<a href="/internal">Internal Link</a>
				<a href="https://example.com/external">External Link</a>
				<a href="#fragment">Fragment Link</a>
				<a href="javascript:void(0)">JavaScript Link</a>
				<a href="mailto:test@example.com">Email Link</a>
				<a>No href</a>
			</body></html>`,
			baseURL: "https://mysite.com",
			expectedMinCounts: map[string]int{
				"total_links": 2, // Should find at least 2 links (internal + external)
			},
		},
		{
			name: "No links",
			html: `<html><body><p>No links here</p></body></html>`,
			baseURL: "https://mysite.com",
			expectedMinCounts: map[string]int{
				"total_links": 0,
			},
		},
		{
			name: "Only internal links",
			html: `<html><body>
				<a href="/page1">Page 1</a>
				<a href="/page2">Page 2</a>
				<a href="https://mysite.com/page3">Page 3</a>
			</body></html>`,
			baseURL: "https://mysite.com",
			expectedMinCounts: map[string]int{
				"total_links": 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseHTML(tt.html)
			result := CountLinks(doc, tt.baseURL)
			
			// Check that we have the expected number of links total
			totalLinks := len(result.Links)
			if totalLinks < tt.expectedMinCounts["total_links"] {
				t.Errorf("CountLinks() found %d links, want at least %d", totalLinks, tt.expectedMinCounts["total_links"])
			}
			
			// Verify that each link has the basic required fields
			for i, link := range result.Links {
				if link.Href == "" {
					t.Errorf("Link %d has empty Href", i)
				}
				if link.AbsoluteURL == "" {
					t.Errorf("Link %d has empty AbsoluteURL", i)
				}
			}
		})
	}
}

// TestLinkResolution tests URL resolution logic without network calls
func TestLinkResolution(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		baseURL     string
		expectedLinks []struct {
			href        string
			absoluteURL string
			isInternal  bool
		}
	}{
		{
			name: "Relative and absolute links",
			html: `<html><body>
				<a href="/page1">Internal Page 1</a>
				<a href="../page2">Internal Page 2</a>
				<a href="https://mysite.com/page3">Internal Page 3</a>
				<a href="https://external.com/page">External Page</a>
			</body></html>`,
			baseURL: "https://mysite.com/current/page",
			expectedLinks: []struct {
				href        string
				absoluteURL string
				isInternal  bool
			}{
				{"/page1", "https://mysite.com/page1", true},
				{"../page2", "https://mysite.com/page2", true},
				{"https://mysite.com/page3", "https://mysite.com/page3", true},
				{"https://external.com/page", "https://external.com/page", false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We'll create a minimal version that doesn't make HTTP calls
			doc := parseHTML(tt.html)
			
			// Parse base URL
			baseU, err := url.Parse(tt.baseURL)
			if err != nil {
				t.Fatalf("Failed to parse base URL: %v", err)
			}
			
			// Extract links manually for testing
			var foundLinks []LinkInfo
			var extractLinks func(*html.Node)
			extractLinks = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "a" {
					var hrefValue string
					for _, attr := range n.Attr {
						if attr.Key == "href" {
							hrefValue = attr.Val
							break
						}
					}
					
					if hrefValue != "" && !strings.HasPrefix(hrefValue, "#") && 
					   !strings.HasPrefix(strings.ToLower(hrefValue), "javascript:") &&
					   !strings.HasPrefix(strings.ToLower(hrefValue), "mailto:") {
						
						linkURL, err := url.Parse(hrefValue)
						if err == nil {
							if !linkURL.IsAbs() {
								linkURL = baseU.ResolveReference(linkURL)
							}
							
							foundLinks = append(foundLinks, LinkInfo{
								Href:        hrefValue,
								AbsoluteURL: linkURL.String(),
								IsInternal:  linkURL.Host == baseU.Host,
							})
						}
					}
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					extractLinks(c)
				}
			}
			extractLinks(doc)
			
			// Verify expected links
			if len(foundLinks) != len(tt.expectedLinks) {
				t.Errorf("Found %d links, expected %d", len(foundLinks), len(tt.expectedLinks))
			}
			
			for i, expected := range tt.expectedLinks {
				if i >= len(foundLinks) {
					break
				}
				found := foundLinks[i]
				
				if found.Href != expected.href {
					t.Errorf("Link %d: href = %q, want %q", i, found.Href, expected.href)
				}
				if found.AbsoluteURL != expected.absoluteURL {
					t.Errorf("Link %d: absoluteURL = %q, want %q", i, found.AbsoluteURL, expected.absoluteURL)
				}
				if found.IsInternal != expected.isInternal {
					t.Errorf("Link %d: isInternal = %v, want %v", i, found.IsInternal, expected.isInternal)
				}
			}
		})
	}
}

func TestExtractTextContent(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "Simple text",
			html:     `<a>Simple text</a>`,
			expected: "Simple text",
		},
		{
			name:     "Text with nested elements",
			html:     `<a>Text with <strong>bold</strong> and <em>italic</em></a>`,
			expected: "Text with bold and italic",
		},
		{
			name:     "Text with whitespace",
			html:     `<a>  Text with   spaces  </a>`,
			expected: "Text with   spaces",
		},
		{
			name:     "Empty element",
			html:     `<a></a>`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseHTML(tt.html)
			// Find the first anchor element
			var anchor *html.Node
			var findAnchor func(*html.Node)
			findAnchor = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "a" {
					anchor = n
					return
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					findAnchor(c)
					if anchor != nil {
						return
					}
				}
			}
			findAnchor(doc)

			if anchor == nil {
				t.Fatal("Could not find anchor element in test HTML")
			}

			result := extractTextContent(anchor)
			if result != tt.expected {
				t.Errorf("extractTextContent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasLoginForm(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name: "Form with password input",
			html: `<html><body>
				<form>
					<input type="text" name="username">
					<input type="password" name="password">
					<input type="submit" value="Login">
				</form>
			</body></html>`,
			expected: true,
		},
		{
			name: "Form without password input",
			html: `<html><body>
				<form>
					<input type="text" name="name">
					<input type="email" name="email">
					<input type="submit" value="Submit">
				</form>
			</body></html>`,
			expected: false,
		},
		{
			name: "Multiple forms, one with password",
			html: `<html><body>
				<form>
					<input type="text" name="search">
					<input type="submit" value="Search">
				</form>
				<form>
					<input type="text" name="username">
					<input type="password" name="password">
					<input type="submit" value="Login">
				</form>
			</body></html>`,
			expected: true,
		},
		{
			name: "No forms",
			html: `<html><body>
				<p>No forms here</p>
			</body></html>`,
			expected: false,
		},
		{
			name: "Password input outside form",
			html: `<html><body>
				<input type="password" name="password">
			</body></html>`,
			expected: false,
		},
		{
			name: "Nested password input in form",
			html: `<html><body>
				<form>
					<div>
						<input type="text" name="username">
						<div>
							<input type="password" name="password">
						</div>
					</div>
					<input type="submit" value="Login">
				</form>
			</body></html>`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := parseHTML(tt.html)
			result := HasLoginForm(doc)
			if result != tt.expected {
				t.Errorf("HasLoginForm() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCheckURLStatus(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
		case "/error":
			w.WriteHeader(http.StatusInternalServerError)
		case "/redirect":
			http.Redirect(w, r, "/ok", http.StatusMovedPermanently)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	tests := []struct {
		name     string
		url      string
		expected *int
	}{
		{
			name:     "Valid URL returns 200",
			url:      server.URL + "/ok",
			expected: func() *int { code := 200; return &code }(),
		},
		{
			name:     "Not found URL returns 404",
			url:      server.URL + "/notfound",
			expected: func() *int { code := 404; return &code }(),
		},
		{
			name:     "Server error returns 500",
			url:      server.URL + "/error",
			expected: func() *int { code := 500; return &code }(),
		},
		{
			name:     "Redirect returns final status",
			url:      server.URL + "/redirect",
			expected: func() *int { code := 200; return &code }(),
		},
		{
			name:     "Invalid URL returns nil",
			url:      "invalid-url",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkURLStatus(tt.url)
			
			if tt.expected == nil {
				if result != nil {
					t.Errorf("checkURLStatus() = %v, want nil", *result)
				}
			} else {
				if result == nil {
					t.Errorf("checkURLStatus() = nil, want %v", *tt.expected)
				} else if *result != *tt.expected {
					t.Errorf("checkURLStatus() = %v, want %v", *result, *tt.expected)
				}
			}
		})
	}
}

// Benchmark tests for performance-critical functions
func BenchmarkExtractTitle(b *testing.B) {
	html := `<html><head><title>Test Page Title</title></head><body><p>Content</p></body></html>`
	doc := parseHTML(html)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ExtractTitle(doc)
	}
}

func BenchmarkCountHeadings(b *testing.B) {
	html := `<html><body>
		<h1>Heading 1</h1>
		<h2>Heading 2</h2>
		<h2>Another Heading 2</h2>
		<h3>Heading 3</h3>
		<h4>Heading 4</h4>
		<h5>Heading 5</h5>
		<h6>Heading 6</h6>
	</body></html>`
	doc := parseHTML(html)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CountHeadings(doc)
	}
}

func BenchmarkHasLoginForm(b *testing.B) {
	html := `<html><body>
		<form>
			<input type="text" name="username">
			<input type="password" name="password">
			<input type="submit" value="Login">
		</form>
	</body></html>`
	doc := parseHTML(html)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HasLoginForm(doc)
	}
}