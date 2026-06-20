#!/bin/bash
# Script to generate GraphQL code using Docker

echo "Generating GraphQL code using Docker..."

docker run --rm \
  -v "$(pwd):/app" \
  -w /app \
  golang:1.24-alpine \
  sh -c "go run github.com/99designs/gqlgen generate"

echo "GraphQL code generation complete!"
