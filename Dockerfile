FROM golang:alpine

RUN apk add --update --no-cache gcc libc-dev

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o server ./cmd/server/main.go

WORKDIR /dist

RUN cp /build/server .

ENV LFLG_JWT_ACCESS_SECRET=access_secret
ENV LFLG_JWT_REFRESH_SECRET=refresh_secret
ENV LFLG_PASS_HASH=$2y$12$w5VUBMSQb3pJKWKDh39gkOyujpd0BscXJY53aETo9oATzXwwi.s02

EXPOSE 8080

CMD ["/dist/server"]
