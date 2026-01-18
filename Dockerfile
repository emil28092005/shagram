FROM golang:1.25.5-alpine AS builder
RUN apk add --no-cache build-base ca-certificates
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 GOOS=linux go build -o shagram ./cmd/server

FROM alpine:latest
RUN apk add --no-cache ca-certificates sqlite-libs sqlite
WORKDIR /app
COPY --from=builder /app/shagram .
COPY static ./static
COPY migrations ./migrations
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["./shagram"]
