# Build stage
FROM golang:1.24.0 AS builder

# Set working directory
WORKDIR /app

# Install librdkafka for Confluent Kafka
RUN apt-get update && \
	apt-get install -y --no-install-recommends \
	build-essential \
	librdkafka-dev \
	pkg-config && \
	rm -rf /var/lib/apt/lists/*

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies with verbose output
RUN go mod download -x

# Tidy up the go.mod and go.sum files
RUN go mod tidy

# Copy source code
COPY . .

# Build all applications using Makefile
RUN go build -o gift-store cmd/app/main.go

# Final stage
FROM ubuntu:22.04

WORKDIR /app

# Install bash, ca-certificates, and librdkafka for Confluent Kafka
RUN apt-get update && \
	apt-get install -y ca-certificates bash librdkafka1 && \
	rm -rf /var/lib/apt/lists/*

# Copy the built applications
COPY --from=builder /app/gift-store .

# Copy configuration files
COPY --from=builder /app/config.yaml .

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./gift-store"]
