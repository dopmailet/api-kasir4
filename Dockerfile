# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code (tanpa tools/)
COPY . .

# Build hanya main package (bukan ./... yang include tools/)
RUN CGO_ENABLED=0 GOOS=linux go build -o kasir-api .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/kasir-api .

# Railway menggunakan PORT environment variable
EXPOSE 8080

CMD ["./kasir-api"]
