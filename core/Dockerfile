FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for downloading dependencies
RUN apk add --no-cache git

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the API binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app/bin/api ./cmd/api/main.go

# Build the Worker binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /app/bin/worker ./cmd/worker/main.go

# Start a new stage from a minimal alpine image
FROM alpine:latest

# Add ca-certificates for HTTPS requests (Nomba API) and tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the pre-built binaries from the previous stage
COPY --from=builder /app/bin/api ./api
COPY --from=builder /app/bin/worker ./worker

# Ensure the binaries are executable
RUN chmod +x ./api ./worker

# Expose port 8080 (default for Kobo API)
EXPOSE 8080

# Command to run the API by default.
# To run the worker, override the command in docker-compose or docker run:
# docker run kobo-core ./worker
CMD ["./api"]
