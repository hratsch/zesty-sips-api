# Use the official Golang image as a parent image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application
CMD ["go", "run", "cmd/server/main.go"]