# Use an official Go image
FROM golang:1.21-alpine

# Install dependencies (including CGO requirements)
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy all files into the container
COPY . .

# Build the application with CGO enabled
RUN CGO_ENABLED=1 go build -o main .

# Expose the app's port (adjust if necessary)
EXPOSE 8080

# Command to run the application
CMD ["./main"]