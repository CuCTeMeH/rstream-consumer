# Dockerfile
# Build stage
FROM golang:1.23.2 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o rstream-consumer .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rstream-consumer .
ENTRYPOINT ["./rstream-consumer"]