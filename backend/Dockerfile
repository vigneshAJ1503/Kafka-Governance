FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Explicitly Tidy mods to embed missing references organically imported above like go-redis
RUN go mod tidy

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o kafka-governance .

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/
COPY --from=builder /app/kafka-governance .

# Create a default .env file in case none is mounted
RUN echo "APP_PORT=8080" > .env

EXPOSE 8080
CMD ["./kafka-governance"]
