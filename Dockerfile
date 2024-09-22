FROM golang:1.21-alpine

# Install git (required by go get)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Verify and tidy up the modules
RUN go mod verify
RUN go mod tidy

# Build the application
RUN go build -o main ./cmd/server

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the application
CMD ["./main"]