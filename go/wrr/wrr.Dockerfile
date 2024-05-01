FROM golang:alpine

# Install iperf3
RUN apk add --no-cache iperf3

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY . .

# Set the GOPATH
ENV GOPATH /go

# Get dependencies using go install
RUN go install -v

# Build the Go app
RUN go build -o wrr .

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./wrr"]