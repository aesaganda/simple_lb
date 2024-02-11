# Start from the Alpine Linux base image
FROM alpine:latest

# Update the package repository and install curl, netcat, and iperf3
RUN apk update && apk add --no-cache curl netcat-openbsd iperf3

# Keep the container running
CMD ["tail", "-f", "/dev/null"]