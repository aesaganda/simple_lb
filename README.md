# Simple Load Balancer

## Description

State of art load balancer for demonstration purposes.

## Installation

### Local

python3 loadBalancer.py --server_host 192.168.1.14 --server_port 8080

### Docker

docker run --name lb --rm -p 8000:8000 -e SERVER_HOST=192.168.1.14 -e SERVER_PORT=8080 loadbalancer

