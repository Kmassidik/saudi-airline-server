# Use the official Go image as a base
FROM golang:1.23 

# Set the Current Working Directory inside the container
WORKDIR /SAUDI-AIRLINE-SERVER

# Copy go.mod and go.sum files into the working directory
COPY go.mod ./
COPY go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Set environment variables (if you have any)
ENV DB_USER=postgres
ENV DB_PASSWORD=admin
ENV DB_NAME=db_saudi_airlines
ENV DB_SSLMODE=disable
ENV DB_HOST=172.17.0.2
ENV URL_IMAGE=http://192.168.1.14:3000/assets/
ENV URL_IMAGE_PROFILE=http://192.168.1.14:3000/images/
ENV JWT_TOKEN=Sunny@day!2024

# Copy the source code and entrypoint script into the container
COPY . .
COPY entrypoint.sh .

# Make the entrypoint script executable
RUN chmod +x entrypoint.sh

# Build the Go app
RUN go build -o my-go-server .

# Use the entrypoint script to run migrations and then the Go app
ENTRYPOINT ["./entrypoint.sh"]
