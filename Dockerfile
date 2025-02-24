# Build
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG APP_VERSION=dev
RUN go build -ldflags="-X 'cherkasyoblenergo_api/internal/config.APP_VERSION=${APP_VERSION}'" -o cherkasyoblenergo_api ./cmd/server/main.go

# Runtime
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/cherkasyoblenergo_api .
COPY .env /app/.env
CMD ["./cherkasyoblenergo_api"]