# Start with a Go runtime image
FROM golang:latest

# Install PostgreSQL client and server
RUN apt-get update && apt-get install -y postgresql postgresql-contrib sudo

# Set up the working directory
WORKDIR /app

# Copy the Go project files to the container
COPY . .

# Install dependencies
RUN go mod download

# Expose the port that the API will listen on
EXPOSE 8080

CMD service postgresql start && sudo -u postgres psql -c "ALTER USER postgres PASSWORD 'goAPI';" && . /app/.env && go run main.go
