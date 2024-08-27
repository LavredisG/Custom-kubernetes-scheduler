# Build stage
FROM golang:alpine as builder

# Install packages used for compiling/building
RUN apk --no-cache add build-base

WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o custom-scheduler .

FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/custom-scheduler /usr/local/bin/custom-scheduler

RUN chmod +x /usr/local/bin/custom-scheduler

# Command to run the scheduler binary
CMD ["/usr/local/bin/custom-scheduler"]

