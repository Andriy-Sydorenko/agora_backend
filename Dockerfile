# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25.2

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates git

# Download deps first (better cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/server

# Final image
FROM alpine:latest
WORKDIR /app

# Install wget for healthcheck
RUN apk add --no-cache wget

COPY --from=builder /bin/app /app/app
COPY config.yml /app/config.yml

EXPOSE 8080
ENTRYPOINT ["/app/app"]
