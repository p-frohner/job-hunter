package professionhu

import (
	"testing"
)

func TestParseJobCards_ReturnsJobs(t *testing.T) {
	html := `<html><body><ul>
		<li data-prof-id="9001"
			data-item-name="Senior Go Developer"
			data-affiliation="Acme Kft.">
			<a href="/allas/senior-go-developer-9001">Job link</a>
			<div class="dsx-card__data-row">Acme Kft.</div>
			<div class="dsx-card__data-row">Budapest</div>
			<div class="dsx-job-card-compact-footer">
				<span class="ds-body-4 ds-color-tertiary">2 napja</span>
			</div>
		</li>
	</ul></body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}

	j := jobs[0]
	if j.ID != "9001" {
		t.Errorf("ID: got %q, want %q", j.ID, "9001")
	}
	if j.Title != "Senior Go Developer" {
		t.Errorf("Title: got %q, want %q", j.Title, "Senior Go Developer")
	}
	if j.Company != "Acme Kft." {
		t.Errorf("Company: got %q, want %q", j.Company, "Acme Kft.")
	}
	if j.Location != "Budapest" {
		t.Errorf("Location: got %q, want %q", j.Location, "Budapest")
	}
	if j.PostedAt != "2 napja" {
		t.Errorf("PostedAt: got %q, want %q", j.PostedAt, "2 napja")
	}
	wantURL := "https://www.profession.hu/allas/senior-go-developer-9001"
	if j.URL != wantURL {
		t.Errorf("URL: got %q, want %q", j.URL, wantURL)
	}
	if j.Source != "professionhu" {
		t.Errorf("Source: got %q, want %q", j.Source, "professionhu")
	}
}

func TestParseJobCards_TitleFallbackToSpan(t *testing.T) {
	// No data-item-name; falls back to DOM selector
	html := `<html><body><ul>
		<li data-prof-id="123">
			<div class="dsx-job-card-compact-header">
				<span class="ds-color-link">
					<a href="/allas/some-job"><span>Fallback Title</span></a>
				</span>
			</div>
			<a href="/allas/some-job">link</a>
		</li>
	</ul></body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Title != "Fallback Title" {
		t.Errorf("Title: got %q, want %q", jobs[0].Title, "Fallback Title")
	}
}

func TestParseJobCards_CompanyFallbackToDataRow(t *testing.T) {
	// No data-affiliation; falls back to first .dsx-card__data-row
	html := `<html><body><ul>
		<li data-prof-id="456" data-item-name="Dev">
			<div class="dsx-card__data-row">Fallback Corp</div>
			<div class="dsx-card__data-row">Debrecen</div>
			<a href="/allas/dev-456">link</a>
		</li>
	</ul></body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Company != "Fallback Corp" {
		t.Errorf("Company: got %q, want %q", jobs[0].Company, "Fallback Corp")
	}
}

func TestParseJobCards_EmptyHTML_ReturnsNil(t *testing.T) {
	jobs, err := parseJobCards("<html><body></body></html>")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestParseJobCards_MissingID_SkipsCard(t *testing.T) {
	// <li> without data-prof-id should not be selected at all (selector is li[data-prof-id])
	// but even if it were, a card with empty ID should be skipped
	html := `<html><body><ul>
		<li>
			<div>No ID here</div>
		</li>
	</ul></body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestParseJobCards_MultipleCards(t *testing.T) {
	html := `<html><body><ul>
		<li data-prof-id="1" data-item-name="Job A"><a href="/allas/a"></a></li>
		<li data-prof-id="2" data-item-name="Job B"><a href="/allas/b"></a></li>
		<li data-prof-id="3" data-item-name="Job C"><a href="/allas/c"></a></li>
	</ul></body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 3 {
		t.Errorf("expected 3 jobs, got %d", len(jobs))
	}
}
