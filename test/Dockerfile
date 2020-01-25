# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /test

# Copy the source from the current directory to the Working Directory inside the container
COPY ./tests ./tests

# Expose port 8000 to the outside world
EXPOSE 4000

# Command to run the executable
CMD ["go test -v ./tests"]
