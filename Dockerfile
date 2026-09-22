FROM golang:1.23-alpine AS build

WORKDIR /src
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -o /out/qaztil ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates \
	&& mkdir -p /app/data
WORKDIR /app
COPY --from=build /out/qaztil /app/qaztil
COPY frontend /app/frontend

ENV WEB_DIR=/app/frontend
ENV DB_PATH=/app/data/qaztil.db

EXPOSE 8080
CMD ["/app/qaztil"]
