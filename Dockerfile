FROM golang:alpine

ARG jwt_access_secret
ARG jwt_refresh_secret
ARG pass_hash

RUN apk add --update --no-cache build-base

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o server ./cmd/server/main.go

WORKDIR /dist

RUN cp /build/server .
RUN mkdir db

ENV LFLG_JWT_ACCESS_SECRET=$jwt_access_secret
ENV LFLG_JWT_REFRESH_SECRET=$jwt_refresh_secret
ENV LFLG_PASS_HASH=$pass_hash

EXPOSE 8080

CMD ["/dist/server"]
