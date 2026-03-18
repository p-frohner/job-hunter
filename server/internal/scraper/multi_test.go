package scraper

import (
	"context"
	"errors"
	"sort"
	"testing"
)

// mockScraper is an in-process Scraper implementation for tests.
type mockScraper struct {
	jobs []Job
	err  error
}

func (m *mockScraper) Search(_ context.Context, _, _ string) ([]Job, error) {
	return m.jobs, m.err
}

func (m *mockScraper) Close() error { return nil }

func TestSearchAll_AggregatesResults(t *testing.T) {
	ms := NewMultiScraper(map[string]Scraper{
		"a": &mockScraper{jobs: []Job{{ID: "1", Title: "Job A"}}},
		"b": &mockScraper{jobs: []Job{{ID: "2", Title: "Job B"}, {ID: "3", Title: "Job C"}}},
	})

	results := ms.SearchAll(context.Background(), "go", "")

	if len(results) != 2 {
		t.Fatalf("expected 2 source results, got %d", len(results))
	}

	totalJobs := 0
	for _, r := range results {
		totalJobs += len(r.Jobs)
		if r.Err != nil {
			t.Errorf("unexpected error for source %q: %v", r.Source, r.Err)
		}
	}
	if totalJobs != 3 {
		t.Errorf("expected 3 total jobs, got %d", totalJobs)
	}
}

func TestSearchAll_OneFailure_OtherSucceeds(t *testing.T) {
	scrapeErr := errors.New("network timeout")
	ms := NewMultiScraper(map[string]Scraper{
		"ok":  &mockScraper{jobs: []Job{{ID: "1", Title: "Good Job"}}},
		"bad": &mockScraper{err: scrapeErr},
	})

	results := ms.SearchAll(context.Background(), "go", "")

	if len(results) != 2 {
		t.Fatalf("expected 2 source results, got %d", len(results))
	}

	// Sort for deterministic assertion
	sort.Slice(results, func(i, j int) bool { return results[i].Source < results[j].Source })

	if results[0].Source != "bad" {
		t.Errorf("unexpected source order")
	}
	if results[0].Err == nil {
		t.Error("expected error for 'bad' source, got nil")
	}
	if results[1].Source != "ok" {
		t.Errorf("unexpected source order")
	}
	if len(results[1].Jobs) != 1 {
		t.Errorf("expected 1 job from 'ok' source, got %d", len(results[1].Jobs))
	}
}

func TestSearchAll_PassesQueryAndLocation(t *testing.T) {
	var gotQuery, gotLocation string
	spy := &spyScraper{fn: func(q, l string) {
		gotQuery = q
		gotLocation = l
	}}

	ms := NewMultiScraper(map[string]Scraper{"s": spy})
	ms.SearchAll(context.Background(), "golang", "Budapest")

	if gotQuery != "golang" {
		t.Errorf("query: got %q, want %q", gotQuery, "golang")
	}
	if gotLocation != "Budapest" {
		t.Errorf("location: got %q, want %q", gotLocation, "Budapest")
	}
}

func TestSearchStream_SendsResultsAndCloses(t *testing.T) {
	ms := NewMultiScraper(map[string]Scraper{
		"x": &mockScraper{jobs: []Job{{ID: "10"}}},
		"y": &mockScraper{jobs: []Job{{ID: "20"}, {ID: "21"}}},
	})

	ch := ms.SearchStream(context.Background(), "go", "")

	var results []SourceResult
	for r := range ch {
		results = append(results, r)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 results from stream, got %d", len(results))
	}
}

func TestSearchStream_ErrorPropagated(t *testing.T) {
	scrapeErr := errors.New("auth wall")
	ms := NewMultiScraper(map[string]Scraper{
		"failing": &mockScraper{err: scrapeErr},
	})

	ch := ms.SearchStream(context.Background(), "go", "")
	r := <-ch

	if r.Err == nil {
		t.Fatal("expected error in stream result, got nil")
	}
	if r.Err.Error() != scrapeErr.Error() {
		t.Errorf("error: got %q, want %q", r.Err.Error(), scrapeErr.Error())
	}
}

func TestSearchStream_ChannelClosedAfterAll(t *testing.T) {
	ms := NewMultiScraper(map[string]Scraper{
		"a": &mockScraper{},
		"b": &mockScraper{},
	})

	ch := ms.SearchStream(context.Background(), "go", "")

	count := 0
	for range ch {
		count++
	}
	if count != 2 {
		t.Errorf("expected 2 results before channel close, got %d", count)
	}
}

// spyScraper records the query/location it was called with.
type spyScraper struct {
	fn func(query, location string)
}

func (s *spyScraper) Search(_ context.Context, query, location string) ([]Job, error) {
	s.fn(query, location)
	return nil, nil
}

func (s *spyScraper) Close() error { return nil }
