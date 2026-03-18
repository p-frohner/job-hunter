package nofluffjobs

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"service-tracker/internal/scraper"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// searchBaseURL is the NoFluffJobs job listing endpoint.
// Query format: ?criteria=keyword%3D<query>+city%3D<location>&sort=newest
const searchBaseURL = "https://nofluffjobs.com/jobs"

// Scraper implements scraper.Scraper for NoFluffJobs.
type Scraper struct {
	browser *rod.Browser
}

// New creates a NoFluffJobs scraper using the provided shared browser instance.
func New(browser *rod.Browser) *Scraper {
	return &Scraper{browser: browser}
}

// Search fetches NoFluffJobs results for the given query and location.
func (s *Scraper) Search(ctx context.Context, query, location string) ([]scraper.Job, error) {
	searchURL := buildSearchURL(query, location)
	slog.Info("nofluffjobs: navigating", "url", searchURL)

	page, err := s.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("nofluffjobs: browser connection lost — network may be unreachable: %w", err)
	}
	defer page.Close()

	page = page.Context(ctx)

	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("nofluffjobs: failed to load %s — check network or SSL proxy: %w", searchURL, err)
	}

	if err := page.WaitLoad(); err != nil {
		slog.Warn("nofluffjobs: WaitLoad error", "error", err)
	}
	time.Sleep(3 * time.Second) // allow React to render job cards

	info := page.MustInfo()
	slog.Info("nofluffjobs: page loaded", "current_url", info.URL, "title", info.Title)

	pageHTML, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("nofluffjobs: failed to get page HTML: %w", err)
	}

	jobs, err := parseJobCards(pageHTML)
	if err != nil {
		return nil, fmt.Errorf("nofluffjobs: failed to parse results: %w", err)
	}

	slog.Info("nofluffjobs: scrape complete", "jobs_found", len(jobs))

	return jobs, nil
}

// Close is a no-op; the shared browser is closed by the caller.
func (s *Scraper) Close() error { return nil }

// buildSearchURL constructs a NoFluffJobs search URL.
func buildSearchURL(query, location string) string {
	criteria := "keyword%3D" + url.QueryEscape(query)
	if location != "" {
		criteria += "+city%3D" + url.QueryEscape(location)
	}
	return searchBaseURL + "?criteria=" + criteria + "&sort=newest"
}
