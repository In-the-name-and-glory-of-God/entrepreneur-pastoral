# Builder stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/server/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Set the user
USER appuser
EXPOSE 8080
# Add a healthcheck (you will need to implement a /health endpoint in your app)
# HEALTHCHECK --interval=30s --timeout=3s \
#   CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./main"]