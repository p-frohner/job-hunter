.PHONY: help install build-all run-server run-client generate

-include .env
export

# The 'help' target will automatically scan this file and print anything with a double hash (##)
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies for both server and client
	cd server && go mod download
	go install github.com/air-verse/air@latest
	cd client && npm install

generate: ## Run codegen for the server (oapi-codegen) and client (orval)
	cd server && oapi-codegen --config server.cfg.yaml ../openapi.yaml
	cd client && npx orval

run-server: ## Run the Go backend with hot reload (Air)
	cd server && air --build.cmd "go build -o tmp/main ./cmd/api/main.go" --build.bin "./tmp/main"

run-client: ## Run the React frontend client
	cd client && npm run dev

docker-up: ## Start everything via Docker Compose
	docker compose up --build

docker-down: ## Stop all Docker services
	docker compose down
