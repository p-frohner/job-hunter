package api

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"service-tracker/internal/scraper"
)

// mockSearcher implements Searcher for tests.
type mockSearcher struct {
	allResults    []scraper.SourceResult
	streamResults []scraper.SourceResult
}

func (m *mockSearcher) SearchAll(_ context.Context, _, _ string) []scraper.SourceResult {
	return m.allResults
}

func (m *mockSearcher) SearchStream(_ context.Context, _, _ string) <-chan scraper.SourceResult {
	ch := make(chan scraper.SourceResult, len(m.streamResults))
	for _, r := range m.streamResults {
		ch <- r
	}
	close(ch)
	return ch
}

// --- SearchJobs (POST /api/search) ---

func TestSearchJobs_MissingQuery_Returns400(t *testing.T) {
	h := NewHandler(&mockSearcher{})

	resp, err := h.SearchJobs(context.Background(), SearchJobsRequestObject{
		Body: &SearchJobsJSONRequestBody{Query: ""},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := resp.(SearchJobs400JSONResponse); !ok {
		t.Errorf("expected 400 response, got %T", resp)
	}
}

func TestSearchJobs_NilBody_Returns400(t *testing.T) {
	h := NewHandler(&mockSearcher{})

	resp, err := h.SearchJobs(context.Background(), SearchJobsRequestObject{Body: nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := resp.(SearchJobs400JSONResponse); !ok {
		t.Errorf("expected 400 response, got %T", resp)
	}
}

func TestSearchJobs_MapsJobFields(t *testing.T) {
	loc := "Budapest"
	snippet := "Great pay"
	source := "linkedin"

	mock := &mockSearcher{
		allResults: []scraper.SourceResult{
			{
				Source: "linkedin",
				Jobs: []scraper.Job{
					{
						ID:       "42",
						Title:    "Go Engineer",
						Company:  "Acme",
						URL:      "https://example.com/job/42",
						Location: loc,
						Snippet:  snippet,
						Source:   source,
					},
				},
			},
		},
	}
	h := NewHandler(mock)

	q := "go"
	resp, err := h.SearchJobs(context.Background(), SearchJobsRequestObject{
		Body: &SearchJobsJSONRequestBody{Query: q},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ok200, ok := resp.(SearchJobs200JSONResponse)
	if !ok {
		t.Fatalf("expected 200 response, got %T", resp)
	}
	if len(ok200.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(ok200.Sources))
	}

	jobs := ok200.Sources[0].Jobs
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}

	j := jobs[0]
	if j.Id != "42" {
		t.Errorf("Id: got %q, want %q", j.Id, "42")
	}
	if j.Title != "Go Engineer" {
		t.Errorf("Title: got %q, want %q", j.Title, "Go Engineer")
	}
	if j.Location == nil || *j.Location != loc {
		t.Errorf("Location: got %v, want %q", j.Location, loc)
	}
	if j.Snippet == nil || *j.Snippet != snippet {
		t.Errorf("Snippet: got %v, want %q", j.Snippet, snippet)
	}
	if j.Source == nil || *j.Source != source {
		t.Errorf("Source: got %v, want %q", j.Source, source)
	}
}

func TestSearchJobs_EmptyOptionalFields_NotSet(t *testing.T) {
	mock := &mockSearcher{
		allResults: []scraper.SourceResult{
			{
				Source: "test",
				Jobs:   []scraper.Job{{ID: "1", Title: "T", Company: "C", URL: "u"}},
			},
		},
	}
	h := NewHandler(mock)

	resp, _ := h.SearchJobs(context.Background(), SearchJobsRequestObject{
		Body: &SearchJobsJSONRequestBody{Query: "x"},
	})

	ok200 := resp.(SearchJobs200JSONResponse)
	j := ok200.Sources[0].Jobs[0]

	if j.Location != nil {
		t.Errorf("Location should be nil, got %v", j.Location)
	}
	if j.PostedAt != nil {
		t.Errorf("PostedAt should be nil, got %v", j.PostedAt)
	}
	if j.Snippet != nil {
		t.Errorf("Snippet should be nil, got %v", j.Snippet)
	}
}

func TestSearchJobs_SourceError_PropagatedInResponse(t *testing.T) {
	scrapeErr := "auth wall hit"
	// Use a wrapped error via fmt.Errorf equivalent
	mock := &mockSearcher{
		allResults: []scraper.SourceResult{
			{Source: "linkedin", Err: errorString(scrapeErr)},
		},
	}
	h := NewHandler(mock)

	resp, _ := h.SearchJobs(context.Background(), SearchJobsRequestObject{
		Body: &SearchJobsJSONRequestBody{Query: "x"},
	})

	ok200 := resp.(SearchJobs200JSONResponse)
	if len(ok200.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(ok200.Sources))
	}
	if ok200.Sources[0].Error == nil {
		t.Fatal("expected error field, got nil")
	}
	if *ok200.Sources[0].Error != scrapeErr {
		t.Errorf("error: got %q, want %q", *ok200.Sources[0].Error, scrapeErr)
	}
}

// --- SearchJobsStream (GET /api/search/stream) ---

func TestSearchJobsStream_MissingQuery_Returns400(t *testing.T) {
	h := NewHandler(&mockSearcher{})

	req := httptest.NewRequest(http.MethodGet, "/api/search/stream", nil)
	w := httptest.NewRecorder()
	h.SearchJobsStream(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestSearchJobsStream_StreamsOneEventPerSource(t *testing.T) {
	mock := &mockSearcher{
		streamResults: []scraper.SourceResult{
			{Source: "linkedin", Jobs: []scraper.Job{{ID: "1", Title: "Job A", Company: "C", URL: "u"}}},
			{Source: "nofluffjobs", Jobs: []scraper.Job{{ID: "2", Title: "Job B", Company: "C", URL: "u"}}},
		},
	}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/search/stream?query=go", nil)
	w := httptest.NewRecorder()
	h.SearchJobsStream(w, req)

	body := w.Body.String()
	dataLines := extractSSEData(body)

	// 2 source events + 1 done event
	if len(dataLines) != 3 {
		t.Fatalf("expected 3 SSE data lines, got %d:\n%s", len(dataLines), body)
	}
	if dataLines[len(dataLines)-1] != `{"done":true}` {
		t.Errorf("last event should be done, got %q", dataLines[len(dataLines)-1])
	}
}

func TestSearchJobsStream_SendsDoneEvent(t *testing.T) {
	h := NewHandler(&mockSearcher{streamResults: nil})

	req := httptest.NewRequest(http.MethodGet, "/api/search/stream?query=go", nil)
	w := httptest.NewRecorder()
	h.SearchJobsStream(w, req)

	body := w.Body.String()
	dataLines := extractSSEData(body)

	if len(dataLines) != 1 {
		t.Fatalf("expected 1 SSE data line (done), got %d", len(dataLines))
	}
	if dataLines[0] != `{"done":true}` {
		t.Errorf("expected done event, got %q", dataLines[0])
	}
}

func TestSearchJobsStream_ContentTypeIsEventStream(t *testing.T) {
	h := NewHandler(&mockSearcher{})

	req := httptest.NewRequest(http.MethodGet, "/api/search/stream?query=go", nil)
	w := httptest.NewRecorder()
	h.SearchJobsStream(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "text/event-stream" {
		t.Errorf("Content-Type: got %q, want %q", ct, "text/event-stream")
	}
}

func TestSearchJobsStream_JobFieldsMappedCorrectly(t *testing.T) {
	loc := "Remote"
	mock := &mockSearcher{
		streamResults: []scraper.SourceResult{
			{Source: "test", Jobs: []scraper.Job{{ID: "99", Title: "T", Company: "C", URL: "u", Location: loc}}},
		},
	}
	h := NewHandler(mock)

	req := httptest.NewRequest(http.MethodGet, "/api/search/stream?query=go", nil)
	w := httptest.NewRecorder()
	h.SearchJobsStream(w, req)

	dataLines := extractSSEData(w.Body.String())
	// First line is the source event
	var sr SourceResult
	if err := json.Unmarshal([]byte(dataLines[0]), &sr); err != nil {
		t.Fatalf("failed to unmarshal event: %v", err)
	}
	if len(sr.Jobs) != 1 {
		t.Fatalf("expected 1 job in event, got %d", len(sr.Jobs))
	}
	if sr.Jobs[0].Location == nil || *sr.Jobs[0].Location != loc {
		t.Errorf("Location: got %v, want %q", sr.Jobs[0].Location, loc)
	}
}

// --- mapJob ---

func TestMapJob_AllFields(t *testing.T) {
	loc := "Budapest"
	postedAt := "2025-03-01"
	snippet := "Great role"
	source := "linkedin"
	desc := "Full description"

	j := mapJob(scraper.Job{
		ID:          "7",
		Title:       "Engineer",
		Company:     "Corp",
		URL:         "https://example.com",
		Location:    loc,
		PostedAt:    postedAt,
		Snippet:     snippet,
		Source:      source,
		Description: desc,
	})

	if j.Id != "7" {
		t.Errorf("Id: got %q", j.Id)
	}
	if j.Location == nil || *j.Location != loc {
		t.Errorf("Location mismatch")
	}
	if j.PostedAt == nil || *j.PostedAt != postedAt {
		t.Errorf("PostedAt mismatch")
	}
	if j.Snippet == nil || *j.Snippet != snippet {
		t.Errorf("Snippet mismatch")
	}
	if j.Source == nil || *j.Source != source {
		t.Errorf("Source mismatch")
	}
	if j.Description == nil || *j.Description != desc {
		t.Errorf("Description mismatch")
	}
}

func TestMapJob_EmptyOptionals_AreNil(t *testing.T) {
	j := mapJob(scraper.Job{ID: "1", Title: "T", Company: "C", URL: "u"})
	if j.Location != nil {
		t.Errorf("Location should be nil")
	}
	if j.PostedAt != nil {
		t.Errorf("PostedAt should be nil")
	}
	if j.Snippet != nil {
		t.Errorf("Snippet should be nil")
	}
	if j.Source != nil {
		t.Errorf("Source should be nil")
	}
	if j.Description != nil {
		t.Errorf("Description should be nil")
	}
}

// --- helpers ---

// extractSSEData pulls the JSON payload from each "data: ..." line in an SSE response body.
func extractSSEData(body string) []string {
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		if after, ok := strings.CutPrefix(line, "data: "); ok {
			lines = append(lines, after)
		}
	}
	return lines
}

// errorString is a minimal error implementation for tests.
type errorString string

func (e errorString) Error() string { return string(e) }
