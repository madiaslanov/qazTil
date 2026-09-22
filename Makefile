.PHONY: run test build swagger fmt

APP := backend/bin/qaztil
SWAG := github.com/swaggo/swag/cmd/swag@v1.16.4

run:
	cd backend && WEB_DIR=../frontend go run ./cmd/api

test:
	cd backend && go test ./...

build:
	mkdir -p backend/bin
	cd backend && go build -o bin/qaztil ./cmd/api

fmt:
	gofmt -w backend/cmd backend/internal

swagger:
	cd backend && go run $(SWAG) init -g main.go -d cmd/api,internal -o internal/adapter/http/docs --parseInternal --parseDependency --exclude internal/adapter/http/docs
