package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// ExtractHtmlVersion determines the HTML version from the DOCTYPE declaration
func ExtractHtmlVersion(doc *html.Node) string {
	// Look for DOCTYPE declaration
	for n := doc.FirstChild; n != nil; n = n.NextSibling {
		if n.Type == html.DoctypeNode {
			// Parse DOCTYPE to determine HTML version
			doctype := strings.ToLower(n.Data)
			if strings.Contains(doctype, "html5") || doctype == "html" {
				return "HTML5"
			} else if strings.Contains(doctype, "xhtml") {
				return "XHTML"
			} else if strings.Contains(doctype, "html 4") {
				return "HTML 4.01"
			}
		}
	}
	return "HTML5" // Default assumption
}

// ExtractTitle retrieves the content of the <title> tag from the HTML document
func ExtractTitle(doc *html.Node) string {
	var title string
	var findTitle func(*html.Node)
	findTitle = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
				title = strings.TrimSpace(n.FirstChild.Data)
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findTitle(c)
			if title != "" {
				return
			}
		}
	}
	findTitle(doc)
	return title
}

// CountHeadings counts the number of each heading level (h1 to h6) in the HTML document
func CountHeadings(doc *html.Node) map[string]int {
	counts := make(map[string]int)
	var countNodes func(*html.Node)
	countNodes = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "h1", "h2", "h3", "h4", "h5", "h6":
				counts[n.Data]++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			countNodes(c)
		}
	}
	countNodes(doc)
	return counts
}

// LinkInfo represents detailed information about a single link
type LinkInfo struct {
	Href        string `json:"href"`         // Original href attribute value
	AbsoluteURL string `json:"absolute_url"` // Resolved absolute URL
	IsInternal  bool   `json:"is_internal"`  // Whether link is internal to the domain
	AnchorText  string `json:"anchor_text"`  // Text content of the link
	StatusCode  *int   `json:"status_code"`  // HTTP status code (nil if not checked)
}

// LinkAnalysis represents the result of link analysis
type LinkAnalysis struct {
	Counts           map[string]int `json:"counts"`
	Links    []LinkInfo     `json:"links"`
}

// CountLinks analyzes and counts internal, external, and inaccessible links in the HTML document
func CountLinks(doc *html.Node, baseURL string) LinkAnalysis {
	result := LinkAnalysis{
		Counts: map[string]int{
			"internal":     0,
			"external":     0,
			"inaccessible": 0,
		},
		Links:    []LinkInfo{},		
	}

	baseU, err := url.Parse(baseURL)
	if err != nil {
		return result
	}

	var countNodes func(*html.Node)
	countNodes = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var hrefValue string
			var anchorText string
			
			// Extract href attribute
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					hrefValue = attr.Val
					break
				}
			}
			
			// Extract anchor text
			anchorText = extractTextContent(n)
			
			// Skip if no href
			if hrefValue == "" {
				return
			}
			
			// Check for inaccessible links first
			if err := validateLinkURL(hrefValue); err != nil || strings.HasPrefix(hrefValue, "#") {
				if !strings.HasPrefix(hrefValue, "#") {
					result.Counts["inaccessible"]++
					result.Links = append(result.Links, LinkInfo{
						Href:        hrefValue,
						AbsoluteURL: "",
						IsInternal:  true,
						AnchorText:  anchorText,
						StatusCode:  nil,
					})
				}
				return
			}

			linkURL, err := url.Parse(hrefValue)
			if err != nil {
				result.Counts["inaccessible"]++
				result.Links = append(result.Links, LinkInfo{
					Href:        hrefValue,
					AbsoluteURL: "",
					IsInternal:  true,
					AnchorText:  anchorText,
					StatusCode:  nil,
				})
				return
			}

			// Resolve relative URLs
			if !linkURL.IsAbs() {
				linkURL = baseU.ResolveReference(linkURL)
			}

			// Create LinkInfo
			linkInfo := LinkInfo{
				Href:        hrefValue,
				AbsoluteURL: linkURL.String(),
				AnchorText:  anchorText,
			}

			// Check status code for http/https URLs
			if linkURL.Scheme == "http" || linkURL.Scheme == "https" {
				statusCode := checkURLStatus(linkURL.String())
				linkInfo.StatusCode = statusCode
			}

			// Categorize link
			if linkURL.Host == baseU.Host {
				linkInfo.IsInternal = true
				result.Counts["internal"]++
				result.Links = append(result.Links, linkInfo)
			} else {
				linkInfo.IsInternal = false
				result.Counts["external"]++
				result.Links = append(result.Links, linkInfo)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			countNodes(c)
		}
	}
	countNodes(doc)
	return result
}

// extractTextContent extracts the text content from a node and its children
func extractTextContent(n *html.Node) string {
	var text strings.Builder
	var extractText func(*html.Node)
	extractText = func(node *html.Node) {
		if node.Type == html.TextNode {
			text.WriteString(node.Data)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(n)
	return strings.TrimSpace(text.String())
}

// checkURLStatus performs a HEAD request to check the status code of a URL
func checkURLStatus(urlStr string) *int {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Allow up to 5 redirects
			if len(via) >= 5 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Try HEAD request first (faster)
	resp, err := client.Head(urlStr)
	if err != nil {
		// If HEAD fails, try GET request
		resp, err = client.Get(urlStr)
		if err != nil {
			// If both fail, return nil (unknown status)
			return nil
		}
	}
	defer resp.Body.Close()

	return &resp.StatusCode
}// validateLinkURL checks if a URL is valid and accessible
func validateLinkURL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("empty URL")
	}
	
	// Skip fragment-only links
	if strings.HasPrefix(rawURL, "#") {
		return fmt.Errorf("fragment-only URL")
	}
	
	// Skip javascript: and mailto: links
	if strings.HasPrefix(strings.ToLower(rawURL), "javascript:") || 
	   strings.HasPrefix(strings.ToLower(rawURL), "mailto:") {
		return fmt.Errorf("non-http URL scheme")
	}
	
	// Try to parse the URL
	_, err := url.Parse(rawURL)
	return err
}

// HasLoginForm checks if the HTML document contains a login form
func HasLoginForm(doc *html.Node) bool {
	var findLoginForm func(*html.Node) bool
	findLoginForm = func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			// Look for forms with password inputs
			if n.Data == "form" {
				hasPasswordInput := false
				var checkInputs func(*html.Node)
				checkInputs = func(node *html.Node) {
					if node.Type == html.ElementNode && node.Data == "input" {
						for _, attr := range node.Attr {
							if attr.Key == "type" && attr.Val == "password" {
								hasPasswordInput = true
								return
							}
						}
					}
					for c := node.FirstChild; c != nil; c = c.NextSibling {
						checkInputs(c)
						if hasPasswordInput {
							return
						}
					}
				}
				checkInputs(n)
				if hasPasswordInput {
					return true
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if findLoginForm(c) {
				return true
			}
		}
		return false
	}
	return findLoginForm(doc)
}