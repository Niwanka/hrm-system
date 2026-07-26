# ==========================================
# STAGE 1: Build Go Binary (Builder)
# ==========================================
FROM golang:alpine AS builder

# Install build tools & ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy all source files
COPY . .

# Ensure dependencies are tidy & download
RUN go mod download

# Build lightweight static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o hrm-backend .

# ==========================================
# STAGE 2: Lightweight Production Container
# ==========================================
FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/hrm-backend .

# Expose Fiber API port
EXPOSE 8080

# Run binary
CMD ["./hrm-backend"]
