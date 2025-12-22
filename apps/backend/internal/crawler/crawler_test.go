package crawler_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"qurio/apps/backend/internal/crawler"
)

func TestCrawl_Depth(t *testing.T) {
	// Setup a mock server
	// / -> links to /level1
	// /level1 -> links to /level2
	// /level2 -> links to /level3
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		switch r.URL.Path {
		case "/":
			fmt.Fprint(w, `<html><body><a href="/level1">Level 1</a></body></html>`)
		case "/level1":
			fmt.Fprint(w, `<html><body><a href="/level2">Level 2</a></body></html>`)
		case "/level2":
			fmt.Fprint(w, `<html><body><a href="/level3">Level 3</a></body></html>`)
		default:
			fmt.Fprint(w, `<html><body>End</body></html>`)
		}
	}))
	defer server.Close()

	// Test Depth 0 (just root)
	c0, _ := crawler.New(crawler.Config{MaxDepth: 0})
	pages0, _ := c0.Crawl(server.URL)
	assert.Len(t, pages0, 1)
	assert.Equal(t, server.URL, pages0[0].URL)

	// Test Depth 1 (root + level1)
	c1, _ := crawler.New(crawler.Config{MaxDepth: 1})
	pages1, _ := c1.Crawl(server.URL)
	assert.Len(t, pages1, 2)
}

func TestCrawl_Exclusion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if r.URL.Path == "/" {
			fmt.Fprint(w, `<html><body><a href="/allowed">Allowed</a><a href="/excluded">Excluded</a></body></html>`)
		} else {
			fmt.Fprint(w, `<html><body>Page</body></html>`)
		}
	}))
	defer server.Close()

	c, err := crawler.New(crawler.Config{
		MaxDepth:   1,
		Exclusions: []string{"excluded"},
	})
	assert.NoError(t, err)

	pages, _ := c.Crawl(server.URL)
	
	// Should contain Root and Allowed, but NOT Excluded
	assert.Len(t, pages, 2)
	for _, p := range pages {
		assert.NotContains(t, p.URL, "excluded")
	}
}

func TestCrawl_Discovery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			fmt.Fprint(w, `<html><body>Root</body></html>`)
		case "/sitemap.xml":
			url := fmt.Sprintf("http://%s/sitemap-page", r.Host)
			fmt.Fprintf(w, `<urlset><url><loc>%s</loc></url></urlset>`, url)
		case "/llms.txt":
			url := fmt.Sprintf("http://%s/llms-page", r.Host)
			fmt.Fprintf(w, `- [LLMs Page](%s)`, url)
		case "/sitemap-page":
			fmt.Fprint(w, "Sitemap Content")
		case "/llms-page":
			fmt.Fprint(w, "LLMs Content")
		}
	}))
	defer server.Close()

	// Use MaxDepth 0 to ensure we rely on discovery, not link following from Root
	c, _ := crawler.New(crawler.Config{MaxDepth: 0})
	pages, err := c.Crawl(server.URL)
	assert.NoError(t, err)

	foundSitemapPage := false
	foundLLMsPage := false
	
	for _, p := range pages {
		if strings.Contains(p.Content, "Sitemap Content") {
			foundSitemapPage = true
		}
		if strings.Contains(p.Content, "LLMs Content") {
			foundLLMsPage = true
		}
	}
	
	assert.True(t, foundSitemapPage, "Should find page linked in sitemap")
	assert.True(t, foundLLMsPage, "Should find page linked in llms.txt")
}