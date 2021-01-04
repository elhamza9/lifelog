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

ENV LFLG_JWT_ACCESS_SECRET=rf2r42CYGV9Mhn3s
ENV LFLG_JWT_REFRESH_SECRET=4h9TvqVLFr2URS3F
ENV LFLG_PASS_HASH=$2y$12$m2fdHd7JRbpIdSphTw7stuJV7zRETmKDv15fIfgCcTPOswn561sGG
ENV LFLG_DB_HOST=192.168.1.199
ENV LFLG_DB_PORT=5432
ENV LFLG_DB_NAME=lifelog
ENV LFLG_DB_USER=postgres
ENV LFLG_DB_PASS=a6XpBmwAx9qm

EXPOSE 8080

CMD ["/dist/server"]
