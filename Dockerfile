# Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG AppVersion
RUN go build -ldflags="-X 'cherkasyoblenergo-api/internal/config.AppVersion=${AppVersion}'" -o cherkasyoblenergo_api ./cmd/server/main.go

# Runtime
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/cherkasyoblenergo_api .
EXPOSE 9011
CMD ["./cherkasyoblenergo_api"]