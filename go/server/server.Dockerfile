FROM alpine:3.19

# Set the Current Working Directory inside the container
WORKDIR /app

COPY server ./

# Expose port 8080 to the outside world

EXPOSE 8080

# Command to run the executable
CMD ["./server"]