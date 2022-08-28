FROM golang:1.19.0-alpine3.16 AS compile

RUN mkdir -p /sda/build
WORKDIR /sda

ADD . .
RUN go build -o build/sda

FROM alpine:3.16

COPY --from=compile /sda/build/sda /usr/local/bin/sda
CMD sda validate --help
