FROM golang:alpine

# Copy Source files to /build directory & download dependencies
WORKDIR /build
COPY . .
RUN go mod download

# Build Go Project & Copy Binary to /dist directory
RUN go build -o server ./cmd/server/main.go
WORKDIR /dist
RUN cp /build/server .

# JWT secrets
ENV LFLG_JWT_ACCESS_SECRET=
ENV LFLG_JWT_REFRESH_SECRE=
# Password Hash
ENV LFLG_PASS_HASH=
# DB Params
ENV LFLG_DB_HOST=
ENV LFLG_DB_PORT=
ENV LFLG_DB_NAME=
ENV LFLG_DB_USER=
ENV LFLG_DB_PASS=

EXPOSE 8080

CMD ["/dist/server"]
