package nofluffjobs

import (
	"strings"
	"testing"
)

func TestParseJobCards_ReturnsJobs(t *testing.T) {
	html := `<html><body>
		<a class="posting" href="/job/acme-go-engineer-abc123">
			<h3 class="posting-title__position">Go Engineer</h3>
			<span class="posting-title__company">Acme Corp</span>
			<ul class="posting-info__location">
				<li>Warsaw</li>
				<li>Remote</li>
			</ul>
			<span class="salary">15 000 – 22 000 PLN</span>
		</a>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}

	j := jobs[0]
	if j.Title != "Go Engineer" {
		t.Errorf("Title: got %q, want %q", j.Title, "Go Engineer")
	}
	if j.Company != "Acme Corp" {
		t.Errorf("Company: got %q, want %q", j.Company, "Acme Corp")
	}
	if j.Location != "Warsaw, Remote" {
		t.Errorf("Location: got %q, want %q", j.Location, "Warsaw, Remote")
	}
	if j.Snippet != "15 000 – 22 000 PLN" {
		t.Errorf("Snippet: got %q, want %q", j.Snippet, "15 000 – 22 000 PLN")
	}
	// ID is derived from last URL segment (full slug)
	if j.ID != "acme-go-engineer-abc123" {
		t.Errorf("ID: got %q, want %q", j.ID, "acme-go-engineer-abc123")
	}
	if j.Source != "nofluffjobs" {
		t.Errorf("Source: got %q, want %q", j.Source, "nofluffjobs")
	}
}

func TestParseJobCards_FallbackCardSelector(t *testing.T) {
	// Uses the second selector: a[href*='/job/']
	html := `<html><body>
		<a href="/job/backend-dev-xyz">
			<h3>Backend Dev</h3>
		</a>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Title != "Backend Dev" {
		t.Errorf("Title: got %q, want %q", jobs[0].Title, "Backend Dev")
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

func TestParseJobCards_MissingTitleOrURL_SkipsCard(t *testing.T) {
	// A card with no title should be skipped
	html := `<html><body>
		<a class="posting" href="/job/no-title-job">
		</a>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs (card skipped), got %d", len(jobs))
	}
}

func TestParseJobCards_LocationFallbackToSingleEl(t *testing.T) {
	// No <li> elements, falls back to container text
	html := `<html><body>
		<a class="posting" href="/job/test-job-1">
			<h3 class="posting-title__position">Test Job</h3>
			<span class="posting-info__location">Kraków</span>
		</a>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if jobs[0].Location != "Kraków" {
		t.Errorf("Location: got %q, want %q", jobs[0].Location, "Kraków")
	}
}

func TestParseJobCards_AbsoluteURLPassedThrough(t *testing.T) {
	html := `<html><body>
		<a class="posting" href="https://nofluffjobs.com/job/existing-url">
			<h3 class="posting-title__position">Full URL Job</h3>
		</a>
	</body></html>`

	jobs, err := parseJobCards(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
	if !strings.HasPrefix(jobs[0].URL, "https://") {
		t.Errorf("URL should be absolute, got %q", jobs[0].URL)
	}
}
