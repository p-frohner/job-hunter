package linkedin

import (
	"log/slog"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"service-tracker/internal/scraper"
)

var entityURNRegex = regexp.MustCompile(`urn:li:jobPosting:(\d+)`)

// parseJobCards extracts jobs from the LinkedIn search results HTML.
func parseJobCards(html string) ([]scraper.Job, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Match only the inner card div, not the <li> wrapper, to avoid duplicates.
	selector := ".base-card.job-search-card"
	slog.Debug("linkedin: querying cards", "selector", selector)

	cards := doc.Find(selector)
	slog.Debug("linkedin: cards found", "count", cards.Length())

	if cards.Length() == 0 {
		excerpt := html
		if len(excerpt) > 3000 {
			excerpt = excerpt[:3000]
		}
		slog.Warn("linkedin: no cards found — page HTML excerpt", "html", excerpt)
		return nil, nil
	}

	var jobs []scraper.Job
	cards.Each(func(i int, s *goquery.Selection) {
		job, err := parseCard(s)
		if err != nil {
			slog.Warn("linkedin: error parsing card", "index", i, "error", err)
			return
		}
		if job.ID == "" {
			slog.Warn("linkedin: card skipped — no ID", "index", i)
			return
		}
		slog.Debug("linkedin: card parsed", "index", i, "id", job.ID, "title", job.Title)
		jobs = append(jobs, job)
	})

	return jobs, nil
}

func parseCard(s *goquery.Selection) (scraper.Job, error) {
	job := scraper.Job{Source: "linkedin"}

	// Job ID — extract from data-entity-urn attribute on the card itself.
	// e.g. data-entity-urn="urn:li:jobPosting:4369496551"
	if urn, exists := s.Attr("data-entity-urn"); exists {
		if m := entityURNRegex.FindStringSubmatch(urn); len(m) > 1 {
			job.ID = m[1]
		}
	}

	// Job URL — from the full-link anchor; href may be relative.
	if href, exists := s.Find("a.base-card__full-link").Attr("href"); exists {
		job.URL = absoluteURL(cleanURL(href))
	}

	// Title
	if t := strings.TrimSpace(s.Find("h3.base-search-card__title").Text()); t != "" {
		job.Title = t
	}

	// Company
	if t := strings.TrimSpace(s.Find("h4.base-search-card__subtitle").Text()); t != "" {
		job.Company = t
	}

	// Location
	if t := strings.TrimSpace(s.Find(".job-search-card__location").Text()); t != "" {
		job.Location = t
	}

	// Time posted
	if timeEl := s.Find("time"); timeEl.Length() > 0 {
		if dt, exists := timeEl.Attr("datetime"); exists {
			job.PostedAt = dt
		} else {
			job.PostedAt = strings.TrimSpace(timeEl.Text())
		}
	}

	return job, nil
}

// cleanURL strips query params from a LinkedIn job URL.
func cleanURL(raw string) string {
	if idx := strings.Index(raw, "?"); idx != -1 {
		raw = raw[:idx]
	}
	return raw
}

// absoluteURL ensures the URL is absolute.
func absoluteURL(u string) string {
	if strings.HasPrefix(u, "http") {
		return u
	}
	return "https://www.linkedin.com" + u
}
