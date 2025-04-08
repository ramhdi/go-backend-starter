FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build with CGO disabled for static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/server ./cmd/server/main.go

# Use distroless as runtime image
FROM gcr.io/distroless/static:nonroot

# Copy binary from builder
COPY --from=builder /app/bin/server /server
COPY --from=builder /app/config.yaml /home/nonroot/config.yaml

# Use non-root user
USER nonroot:nonroot

# Set working directory to nonroot home
WORKDIR /home/nonroot

# Expose API port
EXPOSE 8080

# Run binary
ENTRYPOINT ["/server"]