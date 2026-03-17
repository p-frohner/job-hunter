package scraper

import "context"

// Job is the normalized output from any scraper.
type Job struct {
	ID          string
	Title       string
	Company     string
	Location    string
	URL         string
	PostedAt    string
	Snippet     string
	Source      string // e.g. "linkedin", "nofluffjobs", "professionhu"
	Description string // full job description from the detail page
}

// Scraper is the extension point for job board implementations.
type Scraper interface {
	Search(ctx context.Context, query, location string) ([]Job, error)
	Close() error
}
