import argparse
import socket

# Create an argument parser
parser = argparse.ArgumentParser(description="A simple load balancer")

# Add arguments for the server host and port
parser.add_argument("--server_host", help="The server host", default="192.168.1.14")
parser.add_argument("--server_port", help="The server port", type=int, default=8000)

# Parse the arguments
args = parser.parse_args()

try:
    # Create a socket to listen for clients
    listen_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    listen_sock.bind(('0.0.0.0', 8000))  # Bind to all interfaces and port 8000
    listen_sock.listen(1)

    # Accept a connection from a client and get the client socket and address
    client_sock, client_addr = listen_sock.accept()

    # Print the client address
    print("Received connection from", client_addr)

    # Create a new socket for forwarding the connection
    server_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

    # Use the host and port from the arguments
    server_host = args.server_host
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

# Handle keyboard interrupt (Ctrl+C)
except KeyboardInterrupt:
    print("\nShutting down...")

# Close the listening socket
finally:
    listen_sock.close()