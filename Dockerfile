# Stage 1: Build the Go application
FROM golang:1.20 as build

WORKDIR /app
COPY . .

# Download dependencies
RUN go mod download

# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./main.go

# Stage 2: Prepare the runtime image
FROM heroku/heroku:20

WORKDIR /app
COPY --from=build /app/app /app/app

# Set environment variables
ENV PORT $PORT
EXPOSE $PORT

# Run the application
CMD ["./app"]
