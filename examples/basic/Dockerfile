# Use the official Golang image as a base
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the module dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o /app/traceapp .

# Expose the service port (optional if your app listens on a specific port)
EXPOSE 8080

# Command to run the app
CMD ["/app/traceapp"]
