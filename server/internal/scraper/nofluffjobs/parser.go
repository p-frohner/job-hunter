package nofluffjobs

import (
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"service-tracker/internal/scraper"
)

// cardSelectors are tried in order until one returns results.
// NoFluffJobs is a React SPA; the DOM structure can vary by A/B test bucket.
var cardSelectors = []string{
	"a.posting",
	"a[href*='/job/']",
	"li.posting-list-item a",
}

// parseJobCards extracts jobs from the NoFluffJobs search results HTML.
// NoFluffJobs is a React SPA; the page must be fully loaded before calling this.
func parseJobCards(html string) ([]scraper.Job, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var cards *goquery.Selection
	usedSelector := ""

	for _, sel := range cardSelectors {
		slog.Debug("nofluffjobs: trying card selector", "selector", sel)
		els := doc.Find(sel)
		if els.Length() > 0 {
			cards = els
			usedSelector = sel
			break
		}
	}

	if cards == nil || cards.Length() == 0 {
		excerpt := html
		if len(excerpt) > 5000 {
			excerpt = excerpt[:5000]
		}
		slog.Warn("nofluffjobs: no cards found — page HTML excerpt", "selector", usedSelector, "html", excerpt)
		return nil, nil
	}

	slog.Debug("nofluffjobs: cards found", "selector", usedSelector, "count", cards.Length())

	var jobs []scraper.Job
	cards.Each(func(i int, s *goquery.Selection) {
		job, err := parseCard(s)
		if err != nil {
			slog.Warn("nofluffjobs: error parsing card", "index", i, "error", err)
			return
		}
		if job.URL == "" || job.Title == "" {
			outerHTML, _ := goquery.OuterHtml(s)
			if len(outerHTML) > 500 {
				outerHTML = outerHTML[:500]
			}
			slog.Warn("nofluffjobs: card skipped — missing URL or title", "index", i, "html", outerHTML)
			return
		}
		slog.Debug("nofluffjobs: card parsed", "index", i, "title", job.Title, "url", job.URL)
		jobs = append(jobs, job)
	})

	return jobs, nil
}

func parseCard(s *goquery.Selection) (scraper.Job, error) {
	job := scraper.Job{Source: "nofluffjobs"}

	// URL — the card itself is the anchor element
	if href, exists := s.Attr("href"); exists {
		job.URL = absoluteURL(href)
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
		if t := strings.TrimSpace(s.Find(sel).Text()); t != "" {
			job.Title = t
			break
		}
	}

	// Company
	for _, sel := range []string{
		".posting-title__company",
		"[class*='company']",
	} {
		if t := strings.TrimSpace(s.Find(sel).Text()); t != "" {
			job.Company = t
			break
		}
	}

	// Location — may list multiple cities
	located := false
	for _, sel := range []string{
		".posting-info__location li",
		".locations li",
		"[class*='location'] li",
	} {
		els := s.Find(sel)
		if els.Length() > 0 {
			locs := make([]string, 0, els.Length())
			els.Each(func(_ int, li *goquery.Selection) {
				if t := strings.TrimSpace(li.Text()); t != "" {
					locs = append(locs, t)
				}
			})
			if len(locs) > 0 {
				job.Location = strings.Join(locs, ", ")
				located = true
				break
			}
		}
	}
	if !located {
		for _, sel := range []string{".posting-info__location", "[class*='location']"} {
			if t := strings.TrimSpace(s.Find(sel).Text()); t != "" {
				job.Location = t
				break
			}
		}
	}

	// Salary / snippet
	for _, sel := range []string{".salary", "[class*='salary']", ".posting-salary"} {
		if t := strings.TrimSpace(s.Find(sel).Text()); t != "" {
			job.Snippet = t
			break
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
