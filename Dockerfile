# Use the official Golang image that's compatible with ARM64 architecture.
# This image will be used for building your application.
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./
# Download Go modules
RUN go mod download

# Copy the go source files
COPY cmd/ cmd/
COPY internal/ internal/
COPY pkg/ pkg/
COPY .env /

# Build the Go app as a static binary.
# This assumes your main function is located in the `cmd/api` directory.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /go/bin/app cmd/api/main.go

# Use a minimal alpine image to keep the final image small
FROM alpine:latest

# Install CA certificates, required for making HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory in the container
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /go/bin/app .

# Expose port 8080 on which your application listens
EXPOSE 8080

# Command to run the binary
CMD ["./app"]
