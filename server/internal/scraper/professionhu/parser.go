package professionhu

import (
	"log/slog"
	"strings"

	"github.com/go-rod/rod"
	"service-tracker/internal/scraper"
)

// parseJobCards extracts jobs from the Profession.hu search results page.
// The site is server-side rendered so all content is available after WaitLoad.
func parseJobCards(page *rod.Page) ([]scraper.Job, error) {
	// Each job card is an <li> element with a data-prof-id attribute.
	selector := "li[data-prof-id]"
	slog.Debug("professionhu: querying cards", "selector", selector)

	cards, err := page.Elements(selector)
	if err != nil {
		return nil, err
	}
	slog.Debug("professionhu: cards found", "count", len(cards))

	if len(cards) == 0 {
		html := page.MustHTML()
		if len(html) > 3000 {
			html = html[:3000]
		}
		slog.Warn("professionhu: no cards found — page HTML excerpt", "html", html)
		return nil, nil
	}

	var jobs []scraper.Job
	for i, card := range cards {
		job, err := parseCard(card)
		if err != nil {
			slog.Warn("professionhu: error parsing card", "index", i, "error", err)
			continue
		}
		if job.ID == "" {
			slog.Warn("professionhu: card skipped — no ID", "index", i)
			continue
		}
		slog.Debug("professionhu: card parsed", "index", i, "id", job.ID, "title", job.Title)
		jobs = append(jobs, job)
	}

	return jobs, nil
}

func parseCard(card *rod.Element) (scraper.Job, error) {
	job := scraper.Job{Source: "professionhu"}

	// Job ID — from data-prof-id attribute
	if id, err := card.Attribute("data-prof-id"); err == nil && id != nil {
		job.ID = *id
	}

	// Title — from data-item-name attribute (fast path) or the anchor text
	if name, err := card.Attribute("data-item-name"); err == nil && name != nil {
		job.Title = strings.TrimSpace(*name)
	}
	if job.Title == "" {
		if el, err := card.Element(".dsx-job-card-compact-header .ds-color-link a span"); err == nil {
			job.Title = strings.TrimSpace(el.MustText())
		}
	}

	// Company — from data-affiliation attribute (fast path) or card data row
	if aff, err := card.Attribute("data-affiliation"); err == nil && aff != nil {
		job.Company = strings.TrimSpace(*aff)
	}
	if job.Company == "" {
		if rows, err := card.Elements(".dsx-card__data-row"); err == nil && len(rows) > 0 {
			job.Company = strings.TrimSpace(rows[0].MustText())
		}
	}

	// Location — second data row, or data-location-id as fallback
	if rows, err := card.Elements(".dsx-card__data-row"); err == nil && len(rows) > 1 {
		job.Location = strings.TrimSpace(rows[1].MustText())
	}

	// URL — anchor pointing to the job detail
	if a, err := card.Element("a[href*='/allas/']"); err == nil {
		if href, err := a.Attribute("href"); err == nil && href != nil {
			job.URL = absoluteURL(*href)
		}
	}

	// Posted date
	if el, err := card.Element(".dsx-job-card-compact-footer .ds-body-4.ds-color-tertiary"); err == nil {
		job.PostedAt = strings.TrimSpace(el.MustText())
	}

	return job, nil
}

// absoluteURL ensures the URL is absolute.
func absoluteURL(u string) string {
	if strings.HasPrefix(u, "http") {
		return u
	}
	return "https://www.profession.hu" + u
}
