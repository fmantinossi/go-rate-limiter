FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal runtime image
FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Run the application
CMD ["./main"] 