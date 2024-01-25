# Use an official Python runtime as a parent image
FROM python:3.7-alpine

# Set the working directory in the container to /app
WORKDIR /app

# Add the current directory contents into the container at /app
ADD . /app

# Install any needed packages specified in requirements.txt
# RUN pip install --no-cache-dir -r requirements.txt

# Make port 8000 available to the world outside this container
EXPOSE 8000

# Define environment variables for server host and port
ENV SERVER_HOST 192.168.1.14
ENV SERVER_PORT 8080

# Run loadBalancer.py when the container launches, passing the server host and port as arguments
CMD python loadBalancer.py --server_host $SERVER_HOST --server_port $SERVER_PORT