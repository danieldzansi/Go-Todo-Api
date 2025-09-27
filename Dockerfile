FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build-time dependencies
RUN apk add --no-cache git

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the sources
COPY . .

# Produce a small, static binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/app .


FROM alpine:3.20

# Install runtime deps (TLS & timezone), create app user
RUN apk add --no-cache ca-certificates tzdata \
	&& addgroup -S app && adduser -S -G app app \
	&& update-ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /out/app /app/app

EXPOSE 8080

# Run as non-root
USER app

# Start the server
ENTRYPOINT ["/app/app"]
