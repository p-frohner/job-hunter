package linkedin

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"service-tracker/internal/scraper"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

const (
	searchBaseURL = "https://www.linkedin.com/jobs/search/"
	authWallPath  = "/authwall"
	checkpointURL = "checkpoint"
)

// Scraper implements scraper.Scraper for LinkedIn.
type Scraper struct {
	browser *rod.Browser
}

// New creates a LinkedIn scraper using the provided shared browser instance.
func New(browser *rod.Browser) *Scraper {
	return &Scraper{browser: browser}
}

// Search fetches the first page of LinkedIn job results for the given query and location.
func (s *Scraper) Search(ctx context.Context, query, location string) ([]scraper.Job, error) {
	searchURL := buildSearchURL(query, location)
	slog.Info("linkedin: navigating", "url", searchURL)

	page, err := s.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("linkedin: browser connection lost — network may be unreachable: %w", err)
	}
	defer page.Close()

	// Apply context timeout to navigation
	page = page.Context(ctx)

	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("linkedin: failed to load %s — check network or SSL proxy: %w", searchURL, err)
	}

	// Wait for page to load
	if err := page.WaitLoad(); err != nil {
		slog.Warn("linkedin: WaitLoad error", "error", err)
	}
	time.Sleep(2 * time.Second) // allow JS to render results

	info := page.MustInfo()
	slog.Info("linkedin: page loaded", "current_url", info.URL, "title", info.Title)

	// Detect auth wall
	if strings.Contains(info.URL, authWallPath) || strings.Contains(info.URL, checkpointURL) {
		return nil, fmt.Errorf("linkedin: hit auth wall — LinkedIn is requiring login")
	}

	jobs, err := parseJobCards(page)
	if err != nil {
		return nil, fmt.Errorf("linkedin: failed to parse results: %w", err)
	}

	slog.Info("linkedin: scrape complete", "jobs_found", len(jobs))

	if len(jobs) == 0 {
		html := page.MustHTML()
		if len(html) > 3000 {
			html = html[:3000]
		}
		slog.Warn("linkedin: no jobs parsed — page HTML excerpt", "html", html)
	}

	return jobs, nil
}

// Close is a no-op; the shared browser is closed by the caller.
func (s *Scraper) Close() error { return nil }

// buildSearchURL constructs a LinkedIn job search URL from query params.
func buildSearchURL(query, location string) string {
	params := url.Values{}
	params.Set("keywords", query)
	params.Set("sortBy", "DD")
	if location != "" {
		params.Set("location", location)
	}
	params.Set("f_TPR", "r86400") // posted in last 24h
	return searchBaseURL + "?" + params.Encode()
}
