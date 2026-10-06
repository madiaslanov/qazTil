.PHONY: run test build swagger fmt docker seed

APP := backend/bin/qaztil
SWAG := github.com/swaggo/swag/cmd/swag@v1.16.4

# JWT_SECRET is required. The fallback is for local runs only; Render generates its own.
run:
	cd backend && JWT_SECRET=$${JWT_SECRET:-dev-secret-change-me-0123456789abcdef} go run .

test:
	cd backend && go test ./...

build:
	mkdir -p backend/bin
	cd backend && go build -tags netgo -ldflags '-s -w' -o bin/qaztil .

fmt:
	gofmt -w backend/main.go backend/internal

swagger:
	cd backend && go run $(SWAG) init -g main.go -d .,internal -o internal/adapter/http/docs --parseInternal --parseDependency --exclude internal/adapter/http/docs

docker:
	docker compose up --build

seed:
	cd backend && go run ./cmd/seed
