# Stage 1: Build the Go application
FROM golang:1.20 as build

WORKDIR /app
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN go build -o app ./cmd/server

# Stage 2: Prepare the runtime image
FROM heroku/heroku:20

WORKDIR /app
COPY --from=build /app/app /app/app

# Set environment variables
ENV PORT 8080

# Run the application
CMD ["./app"]
