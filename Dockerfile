# Use the latest version of Go as the base image
FROM golang:1.21 AS base

# create a build artifact
FROM base AS builder
# Set the working directory to the root of the project
WORKDIR /app

# Copy the source code
COPY . .

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build ./cmd/pets

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./pets"]