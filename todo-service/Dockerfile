# Use an official Go runtime as a parent image
FROM golang:1.19 AS build

# Set the working directory inside the container
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Build the Go application
RUN go build -o main

# Start a new, lightweight image
FROM debian:bullseye-slim

# Set the working directory for the final image
WORKDIR /app

# Copy the compiled Go binary from the build image
COPY --from=build /app/main .

# Expose the port your GraphQL server will listen on
EXPOSE 8080

# Define the command to run your GraphQL server
CMD ["./main"]
