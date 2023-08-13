# Start with a Go runtime image
FROM golang:latest AS build

# Set up the working directory
WORKDIR /app

# Copy the Go project files to the container
COPY . .

# Change the working directory to where your main Go file is located
WORKDIR /app/cmd/book_keeper

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/book_keeper

# Start a new stage using a minimal base image
FROM alpine:latest

# Copy the built Go application from the previous stage
COPY --from=build /app/bin/book_keeper /app/book_keeper
RUN chmod +x /app/book_keeper

# Install additional packages for Alpine
RUN apk add --no-cache bash

# Download and make wait-for-it.sh executable
ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

# Expose the port that the API will listen on
EXPOSE 8080

# Start the Go application
CMD ["/app/book_keeper"]
