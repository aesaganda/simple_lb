import argparse
import socket
import threading
import itertools

# Create an argument parser
parser = argparse.ArgumentParser(description="A simple load balancer")

# Add arguments for the server host and port
parser.add_argument("--server_hosts", help="The server hosts, comma separated", default="192.168.1.14,192.168.1.15")
parser.add_argument("--server_weights", help="The server weights, comma separated", default="1,2")
parser.add_argument("--server_port", help="The server port", type=int, default=8000)

# Parse the arguments
args = parser.parse_args()

# Split the server hosts and weights into lists
server_hosts = args.server_hosts.split(',')
server_weights = list(map(int, args.server_weights.split(',')))

# Create a list of servers, repeating each server according to its weight
servers = list(itertools.chain.from_iterable(itertools.repeat(server, weight) for server, weight in zip(server_hosts, server_weights)))

# Keep track of the current server
current_server = 0

# Function to handle client connection
def handle_client(client_sock):
    global current_server

    # Create a new socket for forwarding the connection
    server_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Use the host and port from the arguments
    server_host = servers[current_server]
    server_port = args.server_port

    # Connect the server socket to the specified host and port
    server_sock.connect((server_host, server_port))

    # Start an infinite loop to handle communication between the client and server
    while True:
        # Receive data from the client
        data = client_sock.recv(1024)
        
        # If no data is received, break the loop
        if not data:
            break

        # Forward the data to the server
        server_sock.sendall(data)

        # Receive the server's response
        server_sock.recv(1024)

    # Close the client and server sockets
    client_sock.close()
    server_sock.close()

    # Move to the next server
    current_server = (current_server + 1) % len(servers)

try:
    # Create a socket to listen for clients
    listen_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    listen_sock.bind(('0.0.0.0', 8000))  # Bind to all interfaces and port 8000
    listen_sock.listen(1)

    while True:
        # Accept a connection from a client and get the client socket and address
        client_sock, client_addr = listen_sock.accept()

        # Print the client address
        print("Received connection from", client_addr)

        # Start a new thread to handle the client
        client_thread = threading.Thread(target=handle_client, args=(client_sock,))
        client_thread.start()

# Handle keyboard interrupt (Ctrl+C)
except KeyboardInterrupt:
    print("\nShutting down...")

# Close the listening socket
finally:
    listen_sock.close()