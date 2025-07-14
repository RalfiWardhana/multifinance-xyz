# Multi-stage build untuk optimized production image
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies yang diperlukan
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application dengan optimasi
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -extldflags '-static'" \
    -a -installsuffix cgo \
    -o main .

# Final stage: create minimal image
FROM alpine:latest

# Install ca-certificates untuk HTTPS requests dan timezone data
RUN apk --no-cache add ca-certificates tzdata mysql-client && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set timezone
ENV TZ=Asia/Jakarta

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy environment file jika ada
COPY --from=builder /app/.env* ./

# Create uploads directory dan set permissions
RUN mkdir -p uploads/ktp uploads/selfie && \
    chown -R appuser:appgroup /app && \
    chmod +x main

# Switch to non-root user untuk security
USER appuser

# Expose port
EXPOSE 8080

# Health check - update untuk work dengan host MySQL
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]