
# Use the latest stable Go version
FROM golang:1.21-alpine

# Set working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage caching
COPY go.mod go.sum ./

# Download dependencies before copying the entire source code
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o receipt_processor

# Expose port 8080 for the application
EXPOSE 8080

# Run the built executable
CMD ["./receipt_processor"]
