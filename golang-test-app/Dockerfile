# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Build the Go application
RUN go build -o main

# Expose port 8080 for the application
EXPOSE 8080

# Command to run the application
CMD ["./main"]
