FROM golang:alpine3.19

# Install iperf3 and htop
RUN apk add --no-cache iperf3 htop

# Expose the default iperf3 server port
EXPOSE 5201

# Run iperf3 in server mode on container startup
CMD ["iperf3", "-s"]