# Stage 1: The Builder
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod ./
# If you don't have a go.sum yet, just copy go.mod
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o orbital-shield .

# Stage 2: The Production Image
FROM alpine:latest
RUN apk --no-cache add ca-certificates procps
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=builder /app/orbital-shield .
# Create the log directory
RUN mkdir telemetry_logs

# The Healthcheck: Checks if the process is alive every 30s
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD pgrep orbital-shield || exit 1

CMD ["./orbital-shield"]