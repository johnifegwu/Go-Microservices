# Start with the official Golang base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app
EXPOSE 3000

# Copy the Go modules files
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application code
COPY . .
# Copy the file you want to read into the container
COPY ./dat/data.sql /app/data.sql
COPY ./dat/schema.sql /app/schema.sql

# Build the Go application
RUN go build -o main .

# Command to run the executable
CMD ["./main"]

