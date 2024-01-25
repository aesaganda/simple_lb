FROM golang:alpine

# Install iperf3
RUN apk add --no-cache iperf3

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go/random ./

# Expose port 8080 to the outside world

EXPOSE 8080

# Command to run the executable
CMD ["./random"]