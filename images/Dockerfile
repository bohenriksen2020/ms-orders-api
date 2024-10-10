# Start from the official Go image for building the application
FROM golang:1.23.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ms-orders-api .

# Create a minimal final image
FROM alpine:latest

# Install any dependencies, e.g., for handling SSL certificates
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the binary from the builder image to the final image
COPY --from=builder /app/ms-orders-api .

# Expose the necessary port (in this case, port 3000 for your web service)
EXPOSE 3000

# Command to run the application
CMD ["./ms-orders-api"]
