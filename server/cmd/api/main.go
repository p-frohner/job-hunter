package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"service-tracker/internal/api"
	"service-tracker/internal/scraper"
	"service-tracker/internal/scraper/linkedin"
	"service-tracker/internal/scraper/nofluffjobs"
	"service-tracker/internal/scraper/professionhu"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	headless := os.Getenv("BROWSER_HEADLESS") != "false"

	slog.Info("starting browser", "headless", headless)

	browser, err := scraper.NewBrowser(headless)
	if err != nil {
		slog.Error("failed to start browser", "error", err)
		os.Exit(1)
	}
	defer browser.Close()

	ls := linkedin.New(browser)
	nfj := nofluffjobs.New(browser)
	pro := professionhu.New(browser)

	multi := scraper.NewMultiScraper(map[string]scraper.Scraper{
		"linkedin":     ls,
		"nofluffjobs":  nfj,
		"professionhu": pro,
	})
	defer multi.Close()

	handler := api.NewHandler(multi)
	strictHandler := api.NewStrictHandler(handler, nil)

	r := chi.NewRouter()
	r.Use(api.CorsMiddleware)
	r.Get("/api/search/stream", handler.SearchJobsStream)
	api.HandlerFromMux(strictHandler, r)

	addr := fmt.Sprintf(":%s", port)
	slog.Info("server listening", "addr", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}
