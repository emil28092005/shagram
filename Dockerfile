FROM golang:1.25.5-alpine AS builder

RUN apk add --no-cache build-base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o shagram ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app
COPY --from=builder /app/shagram .
COPY static ./static
COPY migrations ./migrations
RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./shagram"]