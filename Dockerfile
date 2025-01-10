# First stage: build the Go binary
FROM golang:1.23.4-alpine AS builder
LABEL authors="manzi"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy

# Copy the entire project
COPY . .

# Build the Go app
RUN go build -o messaging-app ./cmd/web/main.go

# Second stage: create image for production
FROM debian:bullseye-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder image
COPY --from=builder /app/messaging-app /app/messaging-app
COPY --from=builder /app/.env /app/.env
COPY --from=builder /app/ca.pem /app/ca.pem
COPY --from=builder /app/template /app/template

# Make sure the binary is executable
RUN chmod +x /app/messaging-app

# Expose required ports
EXPOSE 3000
EXPOSE 4000

# Define the default command to run the app
CMD ["/app/messaging-app"]
