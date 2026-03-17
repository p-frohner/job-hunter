package scraper

import (
	"context"
	"log/slog"
	"sync"
)

// SourceResult carries the jobs (or an error) from a single scraper.
type SourceResult struct {
	Source string
	Jobs   []Job
	Err    error
}

// MultiScraper runs multiple Scrapers concurrently and aggregates their results.
// A single scraper failure does not fail the whole search.
type MultiScraper struct {
	scrapers map[string]Scraper
}

// NewMultiScraper creates a MultiScraper from a named map of scrapers.
func NewMultiScraper(scrapers map[string]Scraper) *MultiScraper {
	return &MultiScraper{scrapers: scrapers}
}

// SearchStream runs all scrapers in parallel and sends each SourceResult to the returned
// channel as soon as it finishes. The channel is closed when all scrapers are done.
func (m *MultiScraper) SearchStream(ctx context.Context, query, location string) <-chan SourceResult {
	ch := make(chan SourceResult, len(m.scrapers))
	var wg sync.WaitGroup

	for source, s := range m.scrapers {
		wg.Add(1)
		go func(source string, s Scraper) {
			defer wg.Done()
			jobs, err := s.Search(ctx, query, location)
			if err != nil {
				slog.Warn("scraper failed", "source", source, "error", err)
			}
			ch <- SourceResult{Source: source, Jobs: jobs, Err: err}
		}(source, s)
	}

	go func() { wg.Wait(); close(ch) }()
	return ch
}

// SearchAll runs all scrapers in parallel and returns one SourceResult per scraper.
func (m *MultiScraper) SearchAll(ctx context.Context, query, location string) []SourceResult {
	results := make([]SourceResult, 0, len(m.scrapers))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for source, s := range m.scrapers {
		wg.Add(1)
		go func(source string, s Scraper) {
			defer wg.Done()
			jobs, err := s.Search(ctx, query, location)
			if err != nil {
				slog.Warn("scraper failed", "source", source, "error", err)
			}
			mu.Lock()
			results = append(results, SourceResult{Source: source, Jobs: jobs, Err: err})
			mu.Unlock()
		}(source, s)
	}

	wg.Wait()
	return results
}

// Close closes all underlying scrapers.
func (m *MultiScraper) Close() error {
	var firstErr error
	for _, s := range m.scrapers {
		if err := s.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
