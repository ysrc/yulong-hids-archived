FROM golang:1.10 as builder

RUN mkdir -p /go/src/yulong-hids
WORKDIR /go/src/yulong-hids
ADD . /go/src/yulong-hids
RUN go build -o ./server/server --ldflags='-w -s -linkmode external -extldflags "-static"' ./server/server.go

FROM alpine
MAINTAINER Jason Cooper "mrderek@protonmail.com"
COPY --from=builder /go/src/yulong-hids/server/server .
COPY --from=builder /go/src/yulong-hids/server/.dockerstart.sh /start.sh
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates && update-ca-certificates
RUN apk add --update tzdata
RUN apk add curl
ENV TZ=Asia/Shanghai
RUN rm -rf /var/cache/apk/*
RUN chmod +x /server /start.sh
ENTRYPOINT /start.sh
