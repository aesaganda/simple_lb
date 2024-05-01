# Start from the latest golang base image
FROM golang:alpine

# Add Maintainer Info
LABEL maintainer="A.Eren Sağanda <erensaganda@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Expose port 8080 to the outside
EXPOSE 8080

# Command to run the executable
CMD ["./server"]