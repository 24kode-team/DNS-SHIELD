BIN      := dns-shield
CMD      := ./cmd/shield
BUILD    := go build -ldflags="-s -w" -o bin/$(BIN) $(CMD)

.PHONY: all build run test lint tidy docker docker-run clean

all: build

## Build binary
build:
	@mkdir -p bin
	$(BUILD)

## Run locally (requires port 53 access — use sudo or see README)
run: build
	./bin/$(BIN) --config configs/shield.yaml

## Run tests
test:
	go test ./... -v -race -count=1

## Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

## Tidy dependencies
tidy:
	go mod tidy

## Build Docker image
docker:
	docker build -t dns-shield:latest .

## Run via Docker Compose
docker-run:
	cd deploy && docker compose up --build

## Clean build artifacts
clean:
	rm -rf bin/

## Download all Go dependencies
deps:
	go mod download

## Show blocklist stats via API
stats:
	curl -s http://localhost:8080/metrics | python3 -m json.tool

## Health check
health:
	curl -s http://localhost:8080/health
