# First Stage
FROM golang:alpine AS builder

# Copy Source files to /build directory & download dependencies
WORKDIR /build
COPY . .
RUN go mod download

# Build Go Project & Copy Binary to /dist directory
RUN go build -o server ./cmd/server/main.go
WORKDIR /dist
RUN cp /build/server .

# Second Stage

FROM builder

EXPOSE 8080

CMD ["/dist/server"]
