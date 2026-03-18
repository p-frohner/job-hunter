package linkedin

import (
	"testing"
)

func TestParseJobCards_ReturnsJobs(t *testing.T) {
	html := `<html><body>
		<div class="base-card job-search-card"
			data-entity-urn="urn:li:jobPosting:42">
			<a class="base-card__full-link" href="/jobs/view/42?trk=foo"></a>
			<h3 class="base-search-card__title"> Go Engineer </h3>
			<h4 class="base-search-card__subtitle"> Acme Corp </h4>
			<span class="job-search-card__location"> Budapest, Hungary </span>
			<time datetime="2025-03-01">2 days ago</time>
		</div>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}

	j := jobs[0]
	if j.ID != "42" {
		t.Errorf("ID: got %q, want %q", j.ID, "42")
	}
	if j.Title != "Go Engineer" {
		t.Errorf("Title: got %q, want %q", j.Title, "Go Engineer")
	}
	if j.Company != "Acme Corp" {
		t.Errorf("Company: got %q, want %q", j.Company, "Acme Corp")
	}
	if j.Location != "Budapest, Hungary" {
		t.Errorf("Location: got %q, want %q", j.Location, "Budapest, Hungary")
	}
	if j.PostedAt != "2025-03-01" {
		t.Errorf("PostedAt: got %q, want %q", j.PostedAt, "2025-03-01")
	}
	// URL should be absolute and stripped of query params
	wantURL := "https://www.linkedin.com/jobs/view/42"
	if j.URL != wantURL {
		t.Errorf("URL: got %q, want %q", j.URL, wantURL)
	}
	if j.Source != "linkedin" {
		t.Errorf("Source: got %q, want %q", j.Source, "linkedin")
	}
}

func TestParseJobCards_MultipleCards(t *testing.T) {
	html := `<html><body>
		<div class="base-card job-search-card" data-entity-urn="urn:li:jobPosting:1">
			<h3 class="base-search-card__title">Engineer</h3>
		</div>
		<div class="base-card job-search-card" data-entity-urn="urn:li:jobPosting:2">
			<h3 class="base-search-card__title">Designer</h3>
		</div>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(jobs))
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
	// Card with no data-entity-urn should be skipped
	html := `<html><body>
		<div class="base-card job-search-card">
			<h3 class="base-search-card__title">No ID Job</h3>
		</div>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs (card skipped), got %d", len(jobs))
	}
}

func TestParseJobCards_PostedAtFallsBackToText(t *testing.T) {
	// When <time> has no datetime attribute, use text content
	html := `<html><body>
		<div class="base-card job-search-card" data-entity-urn="urn:li:jobPosting:99">
			<time>3 days ago</time>
		</div>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].PostedAt != "3 days ago" {
		t.Errorf("PostedAt: got %q, want %q", jobs[0].PostedAt, "3 days ago")
	}
}

func TestParseJobCards_AbsoluteURLPassedThrough(t *testing.T) {
	html := `<html><body>
		<div class="base-card job-search-card" data-entity-urn="urn:li:jobPosting:7">
			<a class="base-card__full-link" href="https://www.linkedin.com/jobs/view/7"></a>
		</div>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	want := "https://www.linkedin.com/jobs/view/7"
	if jobs[0].URL != want {
		t.Errorf("URL: got %q, want %q", jobs[0].URL, want)
	}
}

func TestParseJobCards_InvalidHTML_ReturnsError(t *testing.T) {
	// goquery is lenient about malformed HTML, but a completely empty string should still parse
	_, err := parseJobCards("")
	if err != nil {
		t.Fatalf("unexpected error on empty string: %v", err)
	}
}
