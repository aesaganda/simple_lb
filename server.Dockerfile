# Use an official Ubuntu runtime as a parent image
FROM ubuntu:latest

# Install iperf3
RUN apt-get update && apt-get install -y iperf3 htop

# Expose the default iperf3 server port
EXPOSE 5201

# Run iperf3 in server mode on container startup
CMD ["iperf3", "-s", "-D"]