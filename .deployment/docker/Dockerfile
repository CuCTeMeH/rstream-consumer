# Dockerfile
# Build stage
FROM golang:1.23.2 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rstream-consumer .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/rstream-consumer .
ENTRYPOINT ["/app/rstream-consumer", "start"]