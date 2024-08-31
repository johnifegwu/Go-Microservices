# Build stage
FROM golang:alpine AS builder

# Install necessary packages
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Copy the SQL files (Optional if they are required at runtime)
COPY ./infrastructure/data.sql /infrastructure/data.sql
COPY ./infrastructure/schema.sql /infrastructure/schema.sql

# Build the main application
RUN go build -o /go/bin/app .

# Final stage
FROM alpine:latest

# Install certificates
RUN apk --no-cache add ca-certificates

# Copy the compiled Go binary from the builder image
COPY --from=builder /go/bin/app /app

# Set the entrypoint to run the application
ENTRYPOINT ["/app"]

# Optional metadata
LABEL Name=gomicroservices Version=0.0.1

# Expose the application port
EXPOSE 3000
