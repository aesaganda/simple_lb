# Start from the Alpine Linux base image
FROM alpine:latest

# Install iperf3
RUN apk add --no-cache iperf3

# Expose the default iperf3 server port
EXPOSE 5201

# Run iperf3 in server mode
CMD ["iperf3", "-s"]