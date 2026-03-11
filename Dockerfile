# Build
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG AppVersion
RUN CGO_ENABLED=1 go build -ldflags="-X 'cherkasyoblenergo-api/internal/config.AppVersion=${AppVersion}'" -o cherkasyoblenergo_api ./cmd/server/main.go

# Runtime
FROM alpine:3.16
RUN apk add --no-cache tzdata
ENV TZ=Europe/Kyiv

WORKDIR /app
COPY --from=builder /app/cherkasyoblenergo_api .
COPY .env /app/.env
RUN mkdir -p /app/db
EXPOSE 8080
CMD ["./cherkasyoblenergo_api"]