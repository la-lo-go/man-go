# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the source code into the container
COPY . .

# Set environment variables for port and IP
ENV API_PORT=7070
ENV API_IP=0.0.0.0

# Build the Go application
RUN go build -o main .

# Expose the port the application listens on
EXPOSE $PORT

# Run the compiled binary
CMD ["./main"]