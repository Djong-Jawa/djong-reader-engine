# Use the official Golang image as the builder
FROM golang:1.24 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project source code, including .env
COPY . .

# Ensure GOOS and GOARCH are set correctly for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./server.go

# Use a minimal base image for production
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Set the working directory in the final container
WORKDIR /root/

# Copy the built binary and .env file from the builder stage
COPY --from=builder /app/server .
COPY --from=builder /app/.env .

# Ensure the binary has execute permissions
RUN chmod +x server

# Expose the application port (update if necessary)
EXPOSE 8089

# Run the GraphQL server
CMD ["./server"]
