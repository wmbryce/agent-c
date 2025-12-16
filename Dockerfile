FROM golang:1.25.4-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o apiserver .

# Build goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.19

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/apiserver .
COPY --from=builder /go/bin/goose .
COPY --from=builder /build/migrations ./migrations

# Expose port
EXPOSE 8080

# Command to run when starting the container.
ENTRYPOINT ["./apiserver"]
