package professionhu

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

// searchBaseURL is the Profession.hu job listing endpoint.
const searchBaseURL = "https://www.profession.hu/allasok"

// Scraper implements scraper.Scraper for Profession.hu.
type Scraper struct {
	browser *rod.Browser
}

// New creates a Profession.hu scraper using the provided shared browser instance.
func New(browser *rod.Browser) *Scraper {
	return &Scraper{browser: browser}
}

// Search fetches Profession.hu results for the given query and location.
func (s *Scraper) Search(ctx context.Context, query, location string) ([]scraper.Job, error) {
	searchURL := buildSearchURL(query, location)
	slog.Info("professionhu: navigating", "url", searchURL)

	page, err := s.browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, fmt.Errorf("professionhu: browser connection lost — network may be unreachable: %w", err)
	}
	defer page.Close()

	page = page.Context(ctx)

	if err := page.Navigate(searchURL); err != nil {
		return nil, fmt.Errorf("professionhu: failed to load %s — check network or SSL proxy: %w", searchURL, err)
	}

	// Profession.hu is server-side rendered; WaitLoad is sufficient.
	if err := page.WaitLoad(); err != nil {
		return nil, fmt.Errorf("professionhu: page failed to load: %w", err)
	}
	time.Sleep(1 * time.Second)

	if info, err := page.Info(); err != nil {
		slog.Warn("professionhu: could not get page info", "error", err)
	} else {
		slog.Info("professionhu: page loaded", "current_url", info.URL, "title", info.Title)
	}

	pageHTML, err := page.HTML()
	if err != nil {
		return nil, fmt.Errorf("professionhu: failed to get page HTML: %w", err)
	}

	jobs, err := parseJobCards(pageHTML)
	if err != nil {
		return nil, fmt.Errorf("professionhu: failed to parse results: %w", err)
	}

	slog.Info("professionhu: scrape complete", "jobs_found", len(jobs))

	return jobs, nil
}

// Close is a no-op; the shared browser is closed by the caller.
func (s *Scraper) Close() error { return nil }

// buildSearchURL constructs a Profession.hu search URL.
func buildSearchURL(query, location string) string {
	params := url.Values{}
	if query != "" {
		params.Set("search", query)
	}
	if location != "" {
		params.Set("location", location)
	}
	return searchBaseURL + "?" + params.Encode()
}
