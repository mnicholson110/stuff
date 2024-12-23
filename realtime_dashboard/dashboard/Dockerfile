# Stage 1: Build the Go binary
FROM golang:latest AS builder

# Set the working directory in the container
WORKDIR /app

# Copy Go source files
COPY main.go go.mod go.sum ./

# Build the Go binary
RUN go build -o app .

# Stage 2: Build the React app
FROM node:20 AS react-builder

# Set the working directory in the container
WORKDIR /my-dashboard

# Copy React project files
COPY my-dashboard/package.json my-dashboard/package-lock.json ./
RUN npm install
COPY my-dashboard/ ./

# Build the React app
RUN npm run build

# Stage 3: Create the final image
FROM debian:bookworm-slim

# Set the working directory in the container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/app .

# Copy the React build from the React builder stage
COPY --from=react-builder /my-dashboard/build ./frontend/build

# Expose port 8081 (the port our Go server will use)
EXPOSE 8081

# Start the Go binary
CMD ["./app"]
