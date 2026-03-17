package nofluffjobs

import (
	"log/slog"
	"strings"

	"github.com/go-rod/rod"
	"service-tracker/internal/scraper"
)

// cardSelectors are tried in order until one returns results.
// NoFluffJobs is a React SPA; the DOM structure can vary by A/B test bucket.
var cardSelectors = []string{
	"a.posting",
	"a[href*='/job/']",
	"li.posting-list-item a",
}

// parseJobCards extracts jobs from the NoFluffJobs search results page.
// NoFluffJobs is a React SPA; the page must be fully loaded before calling this.
func parseJobCards(page *rod.Page) ([]scraper.Job, error) {
	var cards rod.Elements
	usedSelector := ""

	for _, sel := range cardSelectors {
		slog.Debug("nofluffjobs: trying card selector", "selector", sel)
		els, err := page.Elements(sel)
		if err == nil && len(els) > 0 {
			cards = els
			usedSelector = sel
			break
		}
	}

	slog.Debug("nofluffjobs: cards found", "selector", usedSelector, "count", len(cards))

	if len(cards) == 0 {
		html := page.MustHTML()
		if len(html) > 5000 {
			html = html[:5000]
		}
		slog.Warn("nofluffjobs: no cards found — page HTML excerpt", "html", html)
		return nil, nil
	}

	var jobs []scraper.Job
	for i, card := range cards {
		job, err := parseCard(card)
		if err != nil {
			slog.Warn("nofluffjobs: error parsing card", "index", i, "error", err)
			continue
		}
		if job.URL == "" || job.Title == "" {
			outerHTML, _ := card.HTML()
			if len(outerHTML) > 500 {
				outerHTML = outerHTML[:500]
			}
			slog.Warn("nofluffjobs: card skipped — missing URL or title", "index", i, "html", outerHTML)
			continue
		}
		slog.Debug("nofluffjobs: card parsed", "index", i, "title", job.Title, "url", job.URL)
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func parseCard(card *rod.Element) (scraper.Job, error) {
	job := scraper.Job{Source: "nofluffjobs"}

	// URL — the card itself is the anchor element
	if href, err := card.Attribute("href"); err == nil && href != nil {
		job.URL = absoluteURL(*href)
		parts := strings.Split(strings.TrimRight(job.URL, "/"), "/")
		if len(parts) > 0 {
			job.ID = parts[len(parts)-1]
		}
	}

	// Title — try multiple selectors
	for _, sel := range []string{
		".posting-title__position",
		"h3[class*='title']",
		"h3",
	} {
		if el, err := card.Element(sel); err == nil {
			if t := strings.TrimSpace(el.MustText()); t != "" {
				job.Title = t
				break
			}
		}
	}

	// Company
	for _, sel := range []string{
		".posting-title__company",
		"[class*='company']",
	} {
		if el, err := card.Element(sel); err == nil {
			if t := strings.TrimSpace(el.MustText()); t != "" {
				job.Company = t
				break
			}
		}
	}

	// Location — may list multiple cities
	for _, sel := range []string{
		".posting-info__location li",
		".locations li",
		"[class*='location'] li",
	} {
		if els, err := card.Elements(sel); err == nil && len(els) > 0 {
			locs := make([]string, 0, len(els))
			for _, el := range els {
				if t := strings.TrimSpace(el.MustText()); t != "" {
					locs = append(locs, t)
				}
			}
			if len(locs) > 0 {
				job.Location = strings.Join(locs, ", ")
				break
			}
		}
	}
	if job.Location == "" {
		for _, sel := range []string{".posting-info__location", "[class*='location']"} {
			if el, err := card.Element(sel); err == nil {
				if t := strings.TrimSpace(el.MustText()); t != "" {
					job.Location = t
					break
				}
			}
		}
	}

	// Salary / snippet
	for _, sel := range []string{".salary", "[class*='salary']", ".posting-salary"} {
		if el, err := card.Element(sel); err == nil {
			if t := strings.TrimSpace(el.MustText()); t != "" {
				job.Snippet = t
				break
			}
		}
	}

	return job, nil
}

// absoluteURL ensures the URL is absolute.
func absoluteURL(u string) string {
	if strings.HasPrefix(u, "http") {
		return u
	}
	return "https://nofluffjobs.com" + u
}
