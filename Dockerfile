# Start from official Golang base image
FROM golang:1.24-alpine AS builder

# Install build dependencies (needed for CGO/SQLite if we use CGO, but pure Go sqlite driver is better strictly. 
# gorm/sqlite uses mattn/go-sqlite3 which REQUIRES CGO enabled and gcc)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Dependency caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the binary
# CGO_ENABLED=1 is mandatory for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o carigo-api ./cmd/api/main.go

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install certificates for ensuring https works if needed
RUN apk --no-cache add ca-certificates

# Copy from builder
COPY --from=builder /app/carigo-api .
COPY --from=builder /app/web ./web

# Render expects PORT env, but we'll default to 8080
ENV PORT=8080
# SQLite data persistence directory
RUN mkdir -p /data
ENV DB_PATH=/data/carigo.db

# Expose port
EXPOSE 8080

# Run
CMD ["./carigo-api"]
