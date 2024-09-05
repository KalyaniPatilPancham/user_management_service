# Use the official Go image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go Modules files first for caching
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Run tests
RUN go test ./...

# Build the Go app
RUN go build -o user-management-service .

# Expose the port the service runs on
EXPOSE 8080

# Run the executable
CMD ["./user-management-service"]
