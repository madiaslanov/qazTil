FROM golang:1.23-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/qaztil ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates su-exec \
	&& addgroup -S app \
	&& adduser -S -G app -u 10001 app \
	&& mkdir -p /app/data
WORKDIR /app
COPY --from=build /out/qaztil /app/qaztil
COPY docker/entrypoint.sh /entrypoint.sh
RUN chmod 755 /entrypoint.sh

ENV DB_PATH=/app/data/qaztil.db
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
	CMD ["/bin/sh", "-c", "wget -qO- http://127.0.0.1:${PORT:-8080}/health >/dev/null"]

ENTRYPOINT ["/entrypoint.sh"]
