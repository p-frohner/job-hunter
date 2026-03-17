# Job Hunter - Go Backend

This is the Go backend for the job hunter. It scrapes job boards (LinkedIn, NoFluffJobs, Profession.hu) using a headless Chromium browser via go-rod.

## Tech Stack

- Language: Go 1.25+
- Router: Chi
- Browser automation: go-rod (Chromium)
- API: OpenAPI / oapi-codegen

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Port the server listens on |
| `BROWSER_HEADLESS` | `true` | Set to `false` to show the browser window during scraping |
| `CHROME_BIN` | system default | Path to the Chromium binary (auto-set in Docker) |

## Local Development

> If you're using Docker (`make docker-up` from the project root), you can skip this section entirely.

### Prerequisites

- Go (1.25+): [Install Go](https://go.dev/doc/install)
- Chromium: Required for scraping. Install via your package manager (e.g. `brew install --cask chromium` on Mac).
- Air (hot reload):
  ```
  go install github.com/air-verse/air@latest
  ```

### Run

```
make run-server
```

The server starts on `http://localhost:8080` by default.
