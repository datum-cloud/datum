# Use the official Go image as a build stage
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o datum-authorization-webhook ./cmd/datum-authorization-webhook

# Use a minimal image for the final container
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /app/datum-authorization-webhook .
ENTRYPOINT ["/app/datum-authorization-webhook"]
