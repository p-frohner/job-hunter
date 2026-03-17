package linkedin

import (
	"log/slog"
	"regexp"
	"strings"

	"github.com/go-rod/rod"
	"service-tracker/internal/scraper"
)

var entityURNRegex = regexp.MustCompile(`urn:li:jobPosting:(\d+)`)

// parseJobCards extracts jobs from the LinkedIn search results page.
func parseJobCards(page *rod.Page) ([]scraper.Job, error) {
	if err := page.WaitLoad(); err != nil {
		return nil, err
	}

	// Match only the inner card div, not the <li> wrapper, to avoid duplicates.
	selector := ".base-card.job-search-card"
	slog.Debug("linkedin: querying cards", "selector", selector)

	cards, err := page.Elements(selector)
	if err != nil {
		return nil, err
	}
	slog.Debug("linkedin: cards found", "count", len(cards))

	if len(cards) == 0 {
		html := page.MustHTML()
		if len(html) > 3000 {
			html = html[:3000]
		}
		slog.Warn("linkedin: no cards found — page HTML excerpt", "html", html)
		return nil, nil
	}

	var jobs []scraper.Job
	for i, card := range cards {
		job, err := parseCard(card)
		if err != nil {
			slog.Warn("linkedin: error parsing card", "index", i, "error", err)
			continue
		}
		if job.ID == "" {
			slog.Warn("linkedin: card skipped — no ID", "index", i)
			continue
		}
		slog.Debug("linkedin: card parsed", "index", i, "id", job.ID, "title", job.Title)
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func parseCard(card *rod.Element) (scraper.Job, error) {
	job := scraper.Job{Source: "linkedin"}

	// Job ID — extract from data-entity-urn attribute on the card itself.
	// e.g. data-entity-urn="urn:li:jobPosting:4369496551"
	if urn, err := card.Attribute("data-entity-urn"); err == nil && urn != nil {
		if m := entityURNRegex.FindStringSubmatch(*urn); len(m) > 1 {
			job.ID = m[1]
		}
	}

	// Job URL — from the full-link anchor; href may be relative.
	if a, err := card.Element("a.base-card__full-link"); err == nil {
		if href, err := a.Attribute("href"); err == nil && href != nil {
			job.URL = absoluteURL(cleanURL(*href))
		}
	}

	// Title
	if el, err := card.Element("h3.base-search-card__title"); err == nil {
		job.Title = strings.TrimSpace(el.MustText())
	}

	// Company
	if el, err := card.Element("h4.base-search-card__subtitle"); err == nil {
		job.Company = strings.TrimSpace(el.MustText())
	}

	// Location
	if el, err := card.Element(".job-search-card__location"); err == nil {
		job.Location = strings.TrimSpace(el.MustText())
	}

	// Time posted
	if el, err := card.Element("time"); err == nil {
		if dt, err := el.Attribute("datetime"); err == nil && dt != nil {
			job.PostedAt = *dt
		} else {
			job.PostedAt = strings.TrimSpace(el.MustText())
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
