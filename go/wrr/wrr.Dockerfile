FROM alpine:3.19

# Install iperf3
RUN apk add --no-cache iperf3

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go/wrr ./

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./wrr"]