.PHONY: run test build swagger fmt docker

APP := backend/bin/qaztil
SWAG := github.com/swaggo/swag/cmd/swag@v1.16.4

run:
	cd backend && go run .

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
