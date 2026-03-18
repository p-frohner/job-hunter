package professionhu

import (
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"service-tracker/internal/scraper"
)

// parseJobCards extracts jobs from the Profession.hu search results HTML.
// The site is server-side rendered so all content is available after WaitLoad.
func parseJobCards(html string) ([]scraper.Job, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// Each job card is an <li> element with a data-prof-id attribute.
	selector := "li[data-prof-id]"
	slog.Debug("professionhu: querying cards", "selector", selector)

	cards := doc.Find(selector)
	slog.Debug("professionhu: cards found", "count", cards.Length())

	if cards.Length() == 0 {
		excerpt := html
		if len(excerpt) > 3000 {
			excerpt = excerpt[:3000]
		}
		slog.Warn("professionhu: no cards found — page HTML excerpt", "html", excerpt)
		return nil, nil
	}

	var jobs []scraper.Job
	cards.Each(func(i int, s *goquery.Selection) {
		job, err := parseCard(s)
		if err != nil {
			slog.Warn("professionhu: error parsing card", "index", i, "error", err)
			return
		}
		if job.ID == "" {
			slog.Warn("professionhu: card skipped — no ID", "index", i)
			return
		}
		slog.Debug("professionhu: card parsed", "index", i, "id", job.ID, "title", job.Title)
		jobs = append(jobs, job)
	})

	return jobs, nil
}

func parseCard(s *goquery.Selection) (scraper.Job, error) {
	job := scraper.Job{Source: "professionhu"}

	// Job ID — from data-prof-id attribute
	if id, exists := s.Attr("data-prof-id"); exists {
		job.ID = id
	}

	// Title — from data-item-name attribute (fast path) or the anchor text
	if name, exists := s.Attr("data-item-name"); exists {
		job.Title = strings.TrimSpace(name)
	}
	if job.Title == "" {
		if t := strings.TrimSpace(s.Find(".dsx-job-card-compact-header .ds-color-link a span").Text()); t != "" {
			job.Title = t
		}
	}

	// Company — from data-affiliation attribute (fast path) or card data row
	if aff, exists := s.Attr("data-affiliation"); exists {
		job.Company = strings.TrimSpace(aff)
	}
	if job.Company == "" {
		if rows := s.Find(".dsx-card__data-row"); rows.Length() > 0 {
			job.Company = strings.TrimSpace(rows.First().Text())
		}
	}

	// Location — second data row
	if rows := s.Find(".dsx-card__data-row"); rows.Length() > 1 {
		job.Location = strings.TrimSpace(rows.Eq(1).Text())
	}

	// URL — anchor pointing to the job detail
	if href, exists := s.Find("a[href*='/allas/']").Attr("href"); exists {
		job.URL = absoluteURL(href)
	}

	// Posted date
	if t := strings.TrimSpace(s.Find(".dsx-job-card-compact-footer .ds-body-4.ds-color-tertiary").Text()); t != "" {
		job.PostedAt = t
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
