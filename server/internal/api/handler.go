package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"service-tracker/internal/scraper"
)

// SearchHandler implements StrictServerInterface.
type SearchHandler struct {
	multi *scraper.MultiScraper
}

// NewHandler creates a SearchHandler with the given MultiScraper.
func NewHandler(m *scraper.MultiScraper) *SearchHandler {
	return &SearchHandler{multi: m}
}

// SearchJobs handles POST /api/search.
func (h *SearchHandler) SearchJobs(ctx context.Context, req SearchJobsRequestObject) (SearchJobsResponseObject, error) {
	if req.Body == nil || req.Body.Query == "" {
		return SearchJobs400JSONResponse{Message: "query is required"}, nil
	}

	location := ""
	if req.Body.Location != nil {
		location = *req.Body.Location
	}

	sourceResults := h.multi.SearchAll(ctx, req.Body.Query, location)

	apiSources := make([]SourceResult, 0, len(sourceResults))
	for _, sr := range sourceResults {
		apiJobs := make([]Job, 0, len(sr.Jobs))
		for _, r := range sr.Jobs {
			j := Job{
				Id:      r.ID,
				Title:   r.Title,
				Company: r.Company,
				Url:     r.URL,
			}
			if r.Location != "" {
				j.Location = &r.Location
			}
			if r.PostedAt != "" {
				j.PostedAt = &r.PostedAt
			}
			if r.Snippet != "" {
				j.Snippet = &r.Snippet
			}
			if r.Source != "" {
				j.Source = &r.Source
			}
			if r.Description != "" {
				j.Description = &r.Description
			}
			apiJobs = append(apiJobs, j)
		}

		s := SourceResult{Source: sr.Source, Jobs: apiJobs}
		if sr.Err != nil {
			msg := sr.Err.Error()
			s.Error = &msg
		}
		apiSources = append(apiSources, s)
	}

	return SearchJobs200JSONResponse{Sources: apiSources}, nil
}

// SearchJobsStream handles GET /api/search/stream using Server-Sent Events.
// It writes one JSON event per source as each scraper completes, then a final {"done":true} event.
func (h *SearchHandler) SearchJobsStream(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "query is required", http.StatusBadRequest)
		return
	}
	location := r.URL.Query().Get("location")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	ch := h.multi.SearchStream(r.Context(), query, location)

	for sr := range ch {
		apiJobs := make([]Job, 0, len(sr.Jobs))
		for _, j := range sr.Jobs {
			aj := Job{
				Id:      j.ID,
				Title:   j.Title,
				Company: j.Company,
				Url:     j.URL,
			}
			if j.Location != "" {
				aj.Location = &j.Location
			}
			if j.PostedAt != "" {
				aj.PostedAt = &j.PostedAt
			}
			if j.Snippet != "" {
				aj.Snippet = &j.Snippet
			}
			if j.Source != "" {
				aj.Source = &j.Source
			}
			if j.Description != "" {
				aj.Description = &j.Description
			}
			apiJobs = append(apiJobs, aj)
		}

		event := SourceResult{Source: sr.Source, Jobs: apiJobs}
		if sr.Err != nil {
			msg := sr.Err.Error()
			event.Error = &msg
		}

		data, err := json.Marshal(event)
		if err != nil {
			continue
		}
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	fmt.Fprintf(w, "data: {\"done\":true}\n\n")
	flusher.Flush()
}

// CorsMiddleware allows requests from the React dev server.
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
