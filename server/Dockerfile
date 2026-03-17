FROM golang:1.25.5 AS dev

RUN apt-get update && apt-get install -y --no-install-recommends \
    chromium ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
