# ── Multi-stage Dockerfile ─────────────────────────────────────────────────────
# Stage 1 (builder): compile the Go binary
# Stage 2 (runner):  copy only the binary into a minimal image
# This keeps the final image small (~10 MB vs ~800 MB with the full SDK).

# ── Stage 1: Build ─────────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

# Install git (needed by `go mod download` for some dependencies)
RUN apk add --no-cache git

WORKDIR /app

# Cache module downloads as a separate layer.
# This layer is only invalidated when go.mod or go.sum changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code and build the binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o /app/server ./cmd/server

# ── Stage 2: Run ───────────────────────────────────────────────────────────────
FROM alpine:3.20 AS runner

# Add CA certificates so the app can make outbound HTTPS calls if needed.
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy only the compiled binary from the builder stage.
COPY --from=builder /app/server .

# Run as a non-root user for security.
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 3000

ENTRYPOINT ["/app/server"]
