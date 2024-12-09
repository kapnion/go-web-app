# Use the official Go image as the base image
FROM golang:1.23.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Use a smaller base image for the final stage
FROM alpine:latest  

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the templates and xsl directories
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/xsl ./xsl
COPY --from=builder /app/fonts ./fonts

# Expose port 8086
EXPOSE 8086

# Command to run the application
CMD ["./main"]
