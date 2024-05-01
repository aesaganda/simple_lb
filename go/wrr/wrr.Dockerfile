FROM golang:alpine

# Install iperf3
RUN apk add --no-cache iperf3

# Set the Current Working Directory inside the container
WORKDIR /app

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o wrr .

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./wrr"]