package crawler

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type Config struct {
	MaxDepth   int
	Exclusions []string
	UserAgent  string
}

type Page struct {
	URL     string
	Content string
}

type Crawler struct {
	cfg        Config
	client     *http.Client
	exclusions []*regexp.Regexp
}

func New(cfg Config) (*Crawler, error) {
	var regexps []*regexp.Regexp
	for _, ex := range cfg.Exclusions {
		r, err := regexp.Compile(ex)
		if err != nil {
			return nil, fmt.Errorf("invalid exclusion regex %q: %w", ex, err)
		}
		regexps = append(regexps, r)
	}

	if cfg.UserAgent == "" {
		cfg.UserAgent = "QurioBot/1.0"
	}

	return &Crawler{
		cfg:        cfg,
		client:     &http.Client{Timeout: 10 * time.Second},
		exclusions: regexps,
	}, nil
}

func (c *Crawler) Crawl(startURL string) ([]Page, error) {
	// 1. Discovery Phase
	seeds := []string{startURL}
	
	if parsed, err := url.Parse(startURL); err == nil {
		root := fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host)
		
		// Attempt /sitemap.xml
		if smLinks, err := c.fetchSitemap(root + "/sitemap.xml"); err == nil {
			seeds = append(seeds, smLinks...)
		}
		// Attempt /llms.txt
		if llmLinks, err := c.fetchLLMsTxt(root + "/llms.txt"); err == nil {
			seeds = append(seeds, llmLinks...)
		}
	}

	visited := make(map[string]bool)
	var pages []Page
	var mu sync.Mutex

	var crawl func(u string, depth int)
	crawl = func(u string, depth int) {
		if depth > c.cfg.MaxDepth {
			return
		}

		// Normalize URL
		parsedURL, err := url.Parse(u)
		if err != nil {
			return
		}
		u = parsedURL.String()

		mu.Lock()
		if visited[u] {
			mu.Unlock()
			return
		}
		visited[u] = true
		mu.Unlock()

		// Check exclusions
		for _, r := range c.exclusions {
			if r.MatchString(u) {
				return
			}
		}

		// Fetch
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			return
		}
		req.Header.Set("User-Agent", c.cfg.UserAgent)

		resp, err := c.client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return
		}
		
		bodyReader := io.LimitReader(resp.Body, 10*1024*1024) 
		bodyBytes, err := io.ReadAll(bodyReader)
		if err != nil {
			return
		}
		
		bodyContent := string(bodyBytes)

		mu.Lock()
		pages = append(pages, Page{URL: u, Content: bodyContent})
		mu.Unlock()

		if depth < c.cfg.MaxDepth {
			links := extractLinks(u, string(bodyContent))
			for _, link := range links {
				crawl(link, depth+1)
			}
		}
	}

	for _, seed := range seeds {
		crawl(seed, 0)
	}

	return pages, nil
}

func (c *Crawler) fetchSitemap(u string) ([]string, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	type URL struct {
		Loc string `xml:"loc"`
	}
	type URLSet struct {
		URLs []URL `xml:"url"`
	}
	var urlSet URLSet
	
	limitReader := io.LimitReader(resp.Body, 10*1024*1024)
	if err := xml.NewDecoder(limitReader).Decode(&urlSet); err != nil {
		return nil, err
	}
	
	var links []string
	for _, u := range urlSet.URLs {
		if u.Loc != "" {
			links = append(links, u.Loc)
		}
	}
	return links, nil
}

func (c *Crawler) fetchLLMsTxt(u string) ([]string, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}
	
	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 2*1024*1024))
	if err != nil {
		return nil, err
	}
	content := string(bodyBytes)
	
	// Regex for markdown links [text](url)
	re := regexp.MustCompile(`\[.*?\]\((.*?)\)`)
	matches := re.FindAllStringSubmatch(content, -1)
	var links []string
	for _, m := range matches {
		if len(m) > 1 {
			links = append(links, m[1])
		}
	}
	return links, nil
}

func extractLinks(baseStr, htmlContent string) []string {
	var links []string
	baseURL, err := url.Parse(baseStr)
	if err != nil {
		return links
	}

	z := html.NewTokenizer(strings.NewReader(htmlContent))

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			t := z.Token()
			if t.Data == "a" {
				for _, a := range t.Attr {
					if a.Key == "href" {
						linkURL, err := url.Parse(a.Val)
						if err != nil {
							continue
						}
						absURL := baseURL.ResolveReference(linkURL)
						if absURL.Scheme == "http" || absURL.Scheme == "https" {
							links = append(links, absURL.String())
						}
					}
				}
			}
		}
	}
	return links
}