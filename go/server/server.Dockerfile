# Start from the latest golang base image
FROM golang:alpine

# Add Maintainer Info
LABEL maintainer="A.Eren Sağanda <erensaganda@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Set the GOPATH
ENV GOPATH /go

# Get dependencies using go install
RUN go install -v

# Build the Go app
RUN go build -o server .

# Expose port 8080 to the outside
EXPOSE 8080

# Command to run the executable
CMD ["./server"]