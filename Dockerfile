# Use the official Golang image as the builder
FROM golang:1.23-alpine AS builder

# Create non-root user for building
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory inside the container
WORKDIR /app
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Copy the Go module files and download dependencies
COPY --chown=appuser:appgroup --chmod=755 go.mod go.sum ./
RUN go mod download

# Copy the entire project source code, including .env
COPY --chown=appuser:appgroup --chmod=755 . .

# Ensure GOOS and GOARCH are set correctly for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./server.go

# Use a minimal base image for production
FROM alpine:3.23.3
RUN apk --no-cache add ca-certificates && addgroup -S appgroup && adduser -S appuser -G appgroup

# Switch to non-root user before any file operations
USER appuser

# Set the working directory in the final container
WORKDIR /home/appuser

# Copy the built binary and .env file from the builder stage with correct ownership and permissions
COPY --from=builder --chown=appuser:appgroup --chmod=755 /app/server .
COPY --from=builder --chown=appuser:appgroup --chmod=755 /app/.env .

# Expose the application port (update if necessary)
EXPOSE 8089

# Run the GraphQL server
CMD ["./server"]

