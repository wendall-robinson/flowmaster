# Use a minimal base image
FROM alpine:3.14

# Set working directory inside the container
WORKDIR /app

# Copy the pre-built Go binary into the container
COPY web/main /app/main

# Make sure the binary is executable
RUN chmod +x /app/main

# Command to run the binary
CMD ["/app/main"]
