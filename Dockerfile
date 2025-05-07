# Dockerfile

FROM golang:alpine

WORKDIR /app

# Install Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the port
EXPOSE 3000

# Command to run the application
CMD ["./main"]
